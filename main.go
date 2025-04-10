// main.go - TELJES KÓD - Piros inline validációval (Név mező) és blokkolással
package main

import (
	"errors"
	"fmt"
	"image/color" // Színekhez
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas" // Canvas objektumokhoz
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding" // Szükséges a skillsList miatt
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme" // Téma konstansokhoz
	"fyne.io/fyne/v2/widget"

	"resume-builder/data" // Győződj meg róla, hogy a modul neve helyes!
	"resume-builder/io"   // Győződj meg róla, hogy a modul neve helyes!
)

// --- Validációs Függvények ---
func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("A név megadása kötelező!")
	}
	return nil
}

func validateAll(data *data.ResumeData) error {
	nameValue, _ := data.ContactInfo.Name.Get()
	if err := validateName(nameValue); err != nil {
		return err // Visszaadjuk az első hibát
	}
	// Ide jöhetnének további ellenőrzések...
	return nil // Minden rendben
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Önéletrajz Készítő")

	var resumeData *data.ResumeData
	resumeData = data.NewResumeData()

	// ---- Szekciók Widgetjei ----

	// -- 1. Elérhetőség --
	nameEntry := widget.NewEntryWithData(resumeData.ContactInfo.Name)
	nameEntry.SetPlaceHolder("Teljes név")
	emailEntry := widget.NewEntryWithData(resumeData.ContactInfo.Email)
	emailEntry.SetPlaceHolder("Email cím")
	phoneEntry := widget.NewEntryWithData(resumeData.ContactInfo.Phone)
	phoneEntry.SetPlaceHolder("Telefonszám")
	websiteEntry := widget.NewEntryWithData(resumeData.ContactInfo.Website)
	websiteEntry.SetPlaceHolder("Weboldal (opcionális)")
	linkedinEntry := widget.NewEntryWithData(resumeData.ContactInfo.LinkedIn)
	linkedinEntry.SetPlaceHolder("LinkedIn profil URL (opcionális)")

	// Hiba "címke" canvas.Text-ként piros színnel
	nameErrorLabel := canvas.NewText("", color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // Piros
	nameErrorLabel.TextSize = theme.CaptionTextSize()
	nameErrorLabel.Alignment = fyne.TextAlignLeading
	nameErrorLabel.Hide()

	contactForm := container.NewVBox(
		widget.NewLabel("Név:"), nameEntry, nameErrorLabel, // Hiba címke itt
		widget.NewLabel("Email:"), emailEntry,
		widget.NewLabel("Telefon:"), phoneEntry,
		widget.NewLabel("Weboldal:"), websiteEntry,
		widget.NewLabel("LinkedIn:"), linkedinEntry,
	)

	// Validátor hozzárendelése a Név mezőhöz (LOGOLÁSSAL ÉS REFRESH-sel)
	nameEntry.Validator = func(text string) error {
		log.Println("Name validator fut...") // Látjuk, hogy elindul-e?

		err := validateName(text) // Ellenőrzés
		if err != nil {
			log.Println("Name érvénytelen, hibaüzenet megjelenítése:", err.Error()) // Látjuk a hibát?
			nameErrorLabel.Text = err.Error()
			nameErrorLabel.Refresh() // Canvas frissítése
			nameErrorLabel.Show()    // Megjelenítés
			contactForm.Refresh()    // <<<--- ÚJ: Próbáljuk frissíteni a szülő konténert is!
			return err
		}
		// Ha nincs hiba:
		log.Println("Name érvényes, hibaüzenet elrejtése.") // Látjuk ezt?
		nameErrorLabel.Text = ""
		nameErrorLabel.Refresh() // Canvas frissítése
		nameErrorLabel.Hide()    // Elrejtés
		contactForm.Refresh()    // <<<--- ÚJ: Próbáljuk frissíteni a szülő konténert is!
		return nil
	}

	// -- 2. Összefoglaló --
	summaryEntry := widget.NewEntryWithData(resumeData.Summary)
	summaryEntry.MultiLine = true
	summaryEntry.SetPlaceHolder("Írj egy rövid bemutatkozást vagy szakmai összefoglalót...")
	summaryEntry.Wrapping = fyne.TextWrapWord
	summarySection := container.NewVBox(summaryEntry)

	// -- 3. Munkatapasztalat --
	var experienceList *widget.List
	showExperienceFormDialog := func(title string, entryToEdit *data.ExperienceEntry, onSave func(*data.ExperienceEntry), parent fyne.Window) {
		companyEntry := widget.NewEntry()
		positionEntry := widget.NewEntry()
		startDateEntry := widget.NewEntry()
		startDateEntry.SetPlaceHolder("ÉÉÉÉ-HH")
		endDateEntry := widget.NewEntry()
		endDateEntry.SetPlaceHolder("ÉÉÉÉ-HH vagy Jelenleg")
		descriptionEntry := widget.NewMultiLineEntry()
		descriptionEntry.SetMinRowsVisible(3)
		if entryToEdit != nil {
			companyVal, _ := entryToEdit.Company.Get()
			companyEntry.SetText(companyVal)
			positionVal, _ := entryToEdit.Position.Get()
			positionEntry.SetText(positionVal)
			startDateVal, _ := entryToEdit.StartDate.Get()
			startDateEntry.SetText(startDateVal)
			endDateVal, _ := entryToEdit.EndDate.Get()
			endDateEntry.SetText(endDateVal)
			descriptionVal, _ := entryToEdit.Description.Get()
			descriptionEntry.SetText(descriptionVal)
		}
		formItems := []*widget.FormItem{widget.NewFormItem("Cég", companyEntry), widget.NewFormItem("Pozíció", positionEntry), widget.NewFormItem("Kezdés dátuma", startDateEntry), widget.NewFormItem("Befejezés dátuma", endDateEntry), widget.NewFormItem("Leírás", descriptionEntry)}
		var formDialog dialog.Dialog
		formDialog = dialog.NewForm(title, "Mentés", "Mégse", formItems, func(ok bool) {
			if !ok {
				return
			}
			currentEntry := entryToEdit
			if currentEntry == nil {
				currentEntry = data.NewExperienceEntry()
			}
			currentEntry.Company.Set(companyEntry.Text)
			currentEntry.Position.Set(positionEntry.Text)
			currentEntry.StartDate.Set(startDateEntry.Text)
			currentEntry.EndDate.Set(endDateEntry.Text)
			currentEntry.Description.Set(descriptionEntry.Text)
			onSave(currentEntry)
		}, parent)
		formDialog.Resize(fyne.NewSize(500, 400))
		formDialog.Show()
	}
	experienceList = widget.NewList(
		func() int { return resumeData.Experience.Length() },
		func() fyne.CanvasObject { // CreateItem
			editButton := widget.NewButton("Szerkesztés", nil)
			deleteButton := widget.NewButton("Törlés", nil)
			buttonGroup := container.NewHBox(editButton, deleteButton)
			labels := container.NewVBox(widget.NewLabel("Pozíció @ Cég"), widget.NewLabel("Dátumok"), widget.NewLabel("Rövid leírás..."))
			return container.NewBorder(nil, nil, nil, buttonGroup, labels)
		},
		func(id widget.ListItemID, item fyne.CanvasObject) { // UpdateItem
			dataItemInterface, err := resumeData.Experience.GetValue(id)
			if err != nil {
				log.Printf("Hiba exp elem lekérésekor (%d): %v\n", id, err)
				return
			}
			if entry, ok := dataItemInterface.(*data.ExperienceEntry); ok {
				borderContainer := item.(*fyne.Container)
				labelsVBox := borderContainer.Objects[0].(*fyne.Container)
				buttonGroup := borderContainer.Objects[1].(*fyne.Container)
				editButton := buttonGroup.Objects[0].(*widget.Button)
				deleteButton := buttonGroup.Objects[1].(*widget.Button)
				labelPosCompany := labelsVBox.Objects[0].(*widget.Label)
				labelDates := labelsVBox.Objects[1].(*widget.Label)
				labelDesc := labelsVBox.Objects[2].(*widget.Label)
				company, _ := entry.Company.Get()
				position, _ := entry.Position.Get()
				labelPosCompany.SetText(fmt.Sprintf("%s @ %s", position, company))
				startDate, _ := entry.StartDate.Get()
				endDate, _ := entry.EndDate.Get()
				labelDates.SetText(fmt.Sprintf("%s - %s", startDate, endDate))
				desc, _ := entry.Description.Get()
				if len(desc) > 50 {
					labelDesc.SetText(desc[:50] + "...")
				} else {
					labelDesc.SetText(desc)
				}
				editButton.OnTapped = func() {
					showExperienceFormDialog("Munkatapasztalat szerkesztése", entry, func(updatedEntry *data.ExperienceEntry) {
						log.Printf("Exp elem (%s @ %s) frissítve.\n", position, company)
						experienceList.Refresh()
					}, myWindow)
				}
				deleteButton.OnTapped = func() {
					companyName, _ := entry.Company.Get()
					posName, _ := entry.Position.Get()
					dialog.ShowConfirm("Törlés megerősítése", fmt.Sprintf("Biztosan törölni szeretnéd?\n(%s @ %s)", posName, companyName), func(confirm bool) {
						if !confirm {
							return
						}
						currentList, getErr := resumeData.Experience.Get()
						if getErr != nil {
							log.Printf("Hiba lista lekérésekor exp törléshez: %v\n", getErr)
							return
						}
						if id < 0 || id >= len(currentList) {
							log.Printf("Hiba exp törléskor: Érvénytelen index (%d), lista hossz: %d\n", id, len(currentList))
							return
						}
						newList := append(currentList[:id], currentList[id+1:]...)
						setErr := resumeData.Experience.Set(newList)
						if setErr != nil {
							log.Println("Hiba exp lista frissítésekor törlés után:", setErr)
							dialog.ShowError(setErr, myWindow)
							return
						}
						experienceList.Refresh()
						log.Printf("Exp elem (%s @ %s) törölve (index: %d).\n", posName, companyName, id)
					}, myWindow)
				}
			} else {
				log.Printf("Hiba: Exp elem típusa nem *data.ExperienceEntry (%d)\n", id)
			}
		},
	)
	addExperienceButton := widget.NewButton("Új tapasztalat hozzáadása", func() {
		showExperienceFormDialog("Új Munkatapasztalat", nil, func(newEntry *data.ExperienceEntry) {
			err := resumeData.Experience.Append(newEntry)
			if err != nil {
				log.Println("Hiba exp hozzáadásakor:", err)
				dialog.ShowError(err, myWindow)
				return
			}
			experienceList.Refresh()
		}, myWindow)
	})
	experienceSection := container.NewBorder(nil, addExperienceButton, nil, nil, experienceList)

	// -- 4. Tanulmányok --
	var educationList *widget.List
	showEducationFormDialog := func(title string, entryToEdit *data.EducationEntry, onSave func(*data.EducationEntry), parent fyne.Window) {
		institutionEntry := widget.NewEntry()
		degreeEntry := widget.NewEntry()
		startDateEntry := widget.NewEntry()
		startDateEntry.SetPlaceHolder("ÉÉÉÉ-HH")
		endDateEntry := widget.NewEntry()
		endDateEntry.SetPlaceHolder("ÉÉÉÉ-HH vagy Végzés éve")
		descriptionEntry := widget.NewMultiLineEntry()
		descriptionEntry.SetMinRowsVisible(3)
		if entryToEdit != nil {
			instVal, _ := entryToEdit.Institution.Get()
			institutionEntry.SetText(instVal)
			degVal, _ := entryToEdit.Degree.Get()
			degreeEntry.SetText(degVal)
			startVal, _ := entryToEdit.StartDate.Get()
			startDateEntry.SetText(startVal)
			endVal, _ := entryToEdit.EndDate.Get()
			endDateEntry.SetText(endVal)
			descVal, _ := entryToEdit.Description.Get()
			descriptionEntry.SetText(descVal)
		}
		formItems := []*widget.FormItem{widget.NewFormItem("Intézmény", institutionEntry), widget.NewFormItem("Végzettség/Szak", degreeEntry), widget.NewFormItem("Kezdés dátuma", startDateEntry), widget.NewFormItem("Befejezés dátuma", endDateEntry), widget.NewFormItem("Leírás (opc.)", descriptionEntry)}
		var formDialog dialog.Dialog
		formDialog = dialog.NewForm(title, "Mentés", "Mégse", formItems, func(ok bool) {
			if !ok {
				return
			}
			currentEntry := entryToEdit
			if currentEntry == nil {
				currentEntry = data.NewEducationEntry()
			}
			currentEntry.Institution.Set(institutionEntry.Text)
			currentEntry.Degree.Set(degreeEntry.Text)
			currentEntry.StartDate.Set(startDateEntry.Text)
			currentEntry.EndDate.Set(endDateEntry.Text)
			currentEntry.Description.Set(descriptionEntry.Text)
			onSave(currentEntry)
		}, parent)
		formDialog.Resize(fyne.NewSize(500, 400))
		formDialog.Show()
	}
	educationList = widget.NewList(func() int { return resumeData.Education.Length() }, func() fyne.CanvasObject {
		editButton := widget.NewButton("Szerkesztés", nil)
		deleteButton := widget.NewButton("Törlés", nil)
		buttonGroup := container.NewHBox(editButton, deleteButton)
		labels := container.NewVBox(widget.NewLabel("Végzettség @ Intézmény"), widget.NewLabel("Dátumok"))
		return container.NewBorder(nil, nil, nil, buttonGroup, labels)
	}, func(id widget.ListItemID, item fyne.CanvasObject) {
		dataItemInterface, err := resumeData.Education.GetValue(id)
		if err != nil {
			log.Printf("Hiba edu elem lekérésekor (%d): %v\n", id, err)
			return
		}
		if entry, ok := dataItemInterface.(*data.EducationEntry); ok {
			borderContainer := item.(*fyne.Container)
			labelsVBox := borderContainer.Objects[0].(*fyne.Container)
			buttonGroup := borderContainer.Objects[1].(*fyne.Container)
			editButton := buttonGroup.Objects[0].(*widget.Button)
			deleteButton := buttonGroup.Objects[1].(*widget.Button)
			labelDegreeInst := labelsVBox.Objects[0].(*widget.Label)
			labelDates := labelsVBox.Objects[1].(*widget.Label)
			inst, _ := entry.Institution.Get()
			degree, _ := entry.Degree.Get()
			labelDegreeInst.SetText(fmt.Sprintf("%s @ %s", degree, inst))
			startDate, _ := entry.StartDate.Get()
			endDate, _ := entry.EndDate.Get()
			labelDates.SetText(fmt.Sprintf("%s - %s", startDate, endDate))
			editButton.OnTapped = func() {
				showEducationFormDialog("Tanulmány szerkesztése", entry, func(updatedEntry *data.EducationEntry) {
					log.Printf("Edu elem (%s @ %s) frissítve.\n", degree, inst)
					educationList.Refresh()
				}, myWindow)
			}
			deleteButton.OnTapped = func() {
				instName, _ := entry.Institution.Get()
				degName, _ := entry.Degree.Get()
				dialog.ShowConfirm("Törlés megerősítése", fmt.Sprintf("Biztosan törölni szeretnéd?\n(%s @ %s)", degName, instName), func(confirm bool) {
					if !confirm {
						return
					}
					currentList, getErr := resumeData.Education.Get()
					if getErr != nil {
						log.Printf("Hiba lista lekérésekor edu törléshez: %v\n", getErr)
						return
					}
					if id < 0 || id >= len(currentList) {
						log.Printf("Hiba edu törléskor: Érvénytelen index (%d), lista hossz: %d\n", id, len(currentList))
						return
					}
					newList := append(currentList[:id], currentList[id+1:]...)
					setErr := resumeData.Education.Set(newList)
					if setErr != nil {
						log.Println("Hiba edu lista frissítésekor törlés után:", setErr)
						dialog.ShowError(setErr, myWindow)
						return
					}
					educationList.Refresh()
					log.Printf("Edu elem (%s @ %s) törölve.\n", degName, instName)
				}, myWindow)
			}
		} else {
			log.Printf("Hiba: Edu elem típusa nem *data.EducationEntry (%d)\n", id)
		}
	})
	addEducationButton := widget.NewButton("Új tanulmány hozzáadása", func() {
		showEducationFormDialog("Új Tanulmány", nil, func(newEntry *data.EducationEntry) {
			err := resumeData.Education.Append(newEntry)
			if err != nil {
				log.Println("Hiba edu hozzáadásakor:", err)
				dialog.ShowError(err, myWindow)
				return
			}
			educationList.Refresh()
		}, myWindow)
	})
	educationSection := container.NewBorder(nil, addEducationButton, nil, nil, educationList)

	// -- 5. Készségek --
	var skillsList *widget.List
	skillsList = widget.NewListWithData(resumeData.Skills, func() fyne.CanvasObject {
		return container.NewBorder(nil, nil, nil, widget.NewButton("Törlés", nil), widget.NewLabel("Skill Placeholder"))
	}, func(item binding.DataItem, obj fyne.CanvasObject) {
		border := obj.(*fyne.Container)
		label := border.Objects[0].(*widget.Label)
		button := border.Objects[1].(*widget.Button)
		boundString := item.(binding.String)
		skillName, _ := boundString.Get()
		label.SetText(skillName)
		button.OnTapped = func() {
			currentIndex := -1
			listData, _ := resumeData.Skills.Get()
			valueToDelete, _ := boundString.Get()
			for i, val := range listData {
				if val == valueToDelete {
					currentIndex = i
					break
				}
			}
			if currentIndex != -1 {
				newList := append(listData[:currentIndex], listData[currentIndex+1:]...)
				err := resumeData.Skills.Set(newList)
				if err != nil {
					log.Println("Hiba a készség törlésekor:", err)
				} else {
					log.Println("Készség törölve:", valueToDelete)
				}
			} else {
				log.Println("Hiba: Törlendő készség nem található:", valueToDelete)
			}
		}
	})
	newSkillEntry := widget.NewEntry()
	newSkillEntry.SetPlaceHolder("Új készség neve...")
	addSkillButton := widget.NewButton("Hozzáadás", func() {
		skill := newSkillEntry.Text
		if skill != "" {
			err := resumeData.Skills.Append(skill)
			if err != nil {
				log.Println("Hiba készség hozzáadásakor:", err)
				dialog.ShowError(err, myWindow)
				return
			}
			newSkillEntry.SetText("")
			log.Println("Készség hozzáadva:", skill)
		}
	})
	inputArea := container.NewBorder(nil, nil, nil, addSkillButton, newSkillEntry)
	skillsSection := container.NewBorder(inputArea, nil, nil, nil, skillsList)

	// ---- Fülek Létrehozása ----
	tabs := container.NewAppTabs(
		container.NewTabItem("Elérhetőség", contactForm),
		container.NewTabItem("Összefoglaló", summarySection),
		container.NewTabItem("Munkatapasztalat", experienceSection),
		container.NewTabItem("Tanulmányok", educationSection),
		container.NewTabItem("Készségek", skillsSection),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	// ---- Gombok (Mentés, Betöltés, Export PDF) ----
	saveButton := widget.NewButton("Mentés .cvx", func() {
		_ = nameEntry.Validate() // Lefuttatja a név validátorát (inline hiba megjelenik)

		// Validáció a dialógus előtt
		validationErr := validateAll(resumeData)
		if validationErr != nil {
			dialog.ShowError(errors.New("Mentés sikertelen! Hibák vannak az űrlapon (pl. Név kötelező). Kérjük, javítsd a pirossal jelölt hibákat."), myWindow)
			return
		}
		// Ha nincs hiba, dialógus megnyitása
		saveDialog := dialog.NewFileSave(
			func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				if writer == nil {
					log.Println("Mentés megszakítva")
					return
				}
				defer writer.Close()
				filePath := writer.URI().Path()
				log.Println("Mentés ide:", filePath)
				err = io.SaveResume(resumeData, filePath)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				dialog.ShowInformation("Mentés sikeres", "Önéletrajz elmentve ide: "+filePath, myWindow)
			}, myWindow)
		saveDialog.SetFileName("oneletrajz.cvx")
		saveDialog.Show()
	})

	loadButton := widget.NewButton("Betöltés .cvx", func() {
		loadDialog := dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				if reader == nil {
					log.Println("Betöltés megszakítva")
					return
				}
				defer reader.Close()
				filePath := reader.URI().Path()
				log.Println("Betöltés innen:", filePath)
				loadedData, err := io.LoadResume(filePath)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				// Adatok frissítése
				nameVal, _ := loadedData.ContactInfo.Name.Get()
				resumeData.ContactInfo.Name.Set(nameVal)
				emailVal, _ := loadedData.ContactInfo.Email.Get()
				resumeData.ContactInfo.Email.Set(emailVal)
				phoneVal, _ := loadedData.ContactInfo.Phone.Get()
				resumeData.ContactInfo.Phone.Set(phoneVal)
				webVal, _ := loadedData.ContactInfo.Website.Get()
				resumeData.ContactInfo.Website.Set(webVal)
				linkedVal, _ := loadedData.ContactInfo.LinkedIn.Get()
				resumeData.ContactInfo.LinkedIn.Set(linkedVal)
				summaryVal, _ := loadedData.Summary.Get()
				resumeData.Summary.Set(summaryVal)
				// Listák frissítése
				expVal, expGetErr := loadedData.Experience.Get()
				setErrExp := resumeData.Experience.Set(expVal)
				if expGetErr != nil || setErrExp != nil {
					log.Printf("Hiba exp lista frissítésekor: GetErr=%v, SetErr=%v\n", expGetErr, setErrExp)
				}
				experienceList.Refresh()
				eduVal, eduGetErr := loadedData.Education.Get()
				setErrEdu := resumeData.Education.Set(eduVal)
				if eduGetErr != nil || setErrEdu != nil {
					log.Printf("Hiba edu lista frissítésekor: GetErr=%v, SetErr=%v\n", eduGetErr, setErrEdu)
				}
				educationList.Refresh()
				skillsVal, skillsGetErr := loadedData.Skills.Get()
				skillsSetErr := resumeData.Skills.Set(skillsVal)
				if skillsGetErr != nil || skillsSetErr != nil {
					log.Printf("Hiba skills lista frissítésekor: GetErr=%v, SetErr=%v\n", skillsGetErr, skillsSetErr)
				}
				dialog.ShowInformation("Betöltés sikeres", "Önéletrajz betöltve innen: "+filePath, myWindow)
			}, myWindow)
		loadDialog.Show()
	})

	exportPdfButton := widget.NewButton("Exportálás PDF-be", func() {
		// Validáció a dialógus előtt
		validationErr := validateAll(resumeData)
		if validationErr != nil {
			dialog.ShowError(errors.New("Exportálás sikertelen! Hibák vannak az űrlapon (pl. Név kötelező). Kérjük, javítsd a pirossal jelölt hibákat."), myWindow)
			return
		}
		// Ha nincs hiba, dialógus megnyitása
		saveDialog := dialog.NewFileSave(
			func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				if writer == nil {
					log.Println("PDF Export megszakítva")
					return
				}
				_ = writer.Close()
				filePath := writer.URI().Path()
				log.Println("PDF exportálás ide:", filePath)
				err = io.ExportPDF(resumeData, filePath)
				if err != nil {
					dialog.ShowError(err, myWindow)
					return
				}
				dialog.ShowInformation("Exportálás sikeres", "Önéletrajz exportálva PDF-be: "+filePath, myWindow)
			}, myWindow)
		saveDialog.SetFileName("oneletrajz.pdf")
		saveDialog.Show()
	})

	buttonBox := container.NewHBox(saveButton, loadButton, exportPdfButton)

	// ---- Teljes ablak tartalmának összeállítása ----
	content := container.NewBorder(nil, buttonBox, nil, nil, tabs)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(700, 600))
	myWindow.ShowAndRun()
}
