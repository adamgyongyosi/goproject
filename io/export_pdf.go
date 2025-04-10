// io/export_pdf.go
package io

import (
	"fmt"
	"log"
	"strings"

	gofpdf "github.com/go-pdf/fpdf" // Az újabb, karbantartott fork használata

	"resume-builder/data" // Győződj meg róla, hogy a modul neve helyes!
)

const (
	pageWidthPDF      = 210 // A4 szélesség (mm)
	pageHeightPDF     = 297 // A4 magasság (mm)
	leftMarginPDF     = 15
	topMarginPDF      = 15
	rightMarginPDF    = 15
	bottomMarginPDF   = 15
	contentWidthPDF   = pageWidthPDF - leftMarginPDF - rightMarginPDF
	lineHeightPDF     = 5 // Alap sortávolság (mm)
	sectionSpacingPDF = 8 // Szekciók közötti térköz
)

// ExportPDF generál egy PDF dokumentumot UTF-8 kódolással, AddUTF8Font használatával.
func ExportPDF(resume *data.ResumeData, filePath string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(leftMarginPDF, topMarginPDF, rightMarginPDF)
	pdf.SetAutoPageBreak(true, bottomMarginPDF)

	// *** FONTOS: UTF-8 font hozzáadása közvetlenül a TTF-ből ***
	// Feltételezi, hogy a DejaVuSans.ttf fájl létezik ott, ahonnan a program fut.
	// Ha pl. egy 'fonts' almappában van: pdf.AddUTF8Font("DejaVu", "", "fonts/DejaVuSans.ttf")
	// Az üres stringek "" azt jelentik, hogy a sima (Regular) és a félkövér (Bold) stílust
	// próbálja ugyanabból a TTF fájlból betölteni, vagy alapértelmezett emulációt használ.
	// Ha van külön Bold TTF (pl. DejaVuSans-Bold.ttf), akkor azt kellene hozzáadni:
	// pdf.AddUTF8Font("DejaVu", "B", "DejaVuSans-Bold.ttf")
	pdf.AddUTF8Font("DejaVu", "", "DejaVuSans.ttf")  // Csak a Regular stílust adjuk hozzá expliciten
	pdf.AddUTF8Font("DejaVu", "B", "DejaVuSans.ttf") // Próbáljuk a Boldot is ugyanabból (vagy használjunk Bold TTF-et)

	pdf.AddPage()

	// Alapértelmezett font beállítása az egész dokumentumra
	pdf.SetFont("DejaVu", "", 10)

	// --- Tartalom hozzáadása (már nem kell a 'tr' fordító) ---
	renderHeaderPDF_UTF8(pdf, resume)

	renderSectionTitlePDF_UTF8(pdf, "Összefoglaló")
	summary, _ := resume.Summary.Get()
	renderMultiLineTextPDF_UTF8(pdf, summary)
	pdf.Ln(sectionSpacingPDF)

	renderSectionTitlePDF_UTF8(pdf, "Munkatapasztalat")
	expList, _ := resume.Experience.Get()
	if len(expList) > 0 {
		for _, item := range expList {
			if entry, ok := item.(*data.ExperienceEntry); ok {
				renderExperienceEntryPDF_UTF8(pdf, entry)
				pdf.Ln(lineHeightPDF / 2)
			}
		}
	} else {
		renderMultiLineTextPDF_UTF8(pdf, "Nincs megadva munkatapasztalat.")
	}
	pdf.Ln(sectionSpacingPDF)

	renderSectionTitlePDF_UTF8(pdf, "Tanulmányok")
	eduList, _ := resume.Education.Get()
	if len(eduList) > 0 {
		for _, item := range eduList {
			if entry, ok := item.(*data.EducationEntry); ok {
				renderEducationEntryPDF_UTF8(pdf, entry)
				pdf.Ln(lineHeightPDF / 2)
			}
		}
	} else {
		renderMultiLineTextPDF_UTF8(pdf, "Nincsenek megadva tanulmányok.")
	}
	pdf.Ln(sectionSpacingPDF)

	// --- Készségek (ÚJ!) ---
	renderSectionTitlePDF_UTF8(pdf, "Készségek") // Nincs tr
	skillsList, errSkillsGet := resume.Skills.Get()
	if errSkillsGet != nil {
		log.Println("Hiba a készségek lekérésekor PDF exporthoz:", errSkillsGet)
		renderMultiLineTextPDF_UTF8(pdf, "Hiba a készségek betöltésekor.")
	} else if len(skillsList) > 0 {
		// Készségek kiírása vesszővel elválasztva
		skillsText := strings.Join(skillsList, ", ")
		renderMultiLineTextPDF_UTF8(pdf, skillsText) // Nincs tr
	} else {
		renderMultiLineTextPDF_UTF8(pdf, "Nincsenek megadva készségek.") // Nincs tr
	}
	pdf.Ln(sectionSpacingPDF) // Térköz a végén

	// --- PDF Mentése ---
	err := pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Printf("Hiba a PDF mentésekor: %v\n", err)
		// Az AddUTF8Font is adhat "cannot find" hibát, ha a .ttf hiányzik
		if strings.Contains(err.Error(), "cannot find") || strings.Contains(err.Error(), "no such file") {
			return fmt.Errorf("PDF mentési hiba: Nem található a szükséges DejaVuSans.ttf betűtípus fájl. %w", err)
		}
		return err
	}
	if pdf.Error() != nil {
		log.Printf("gofpdf belső hiba: %v\n", pdf.Error())
		return fmt.Errorf("PDF generálási hiba: %w", pdf.Error())
	}
	return nil
}

// --- Segédfüggvények (_UTF8 végződéssel az egyértelműség kedvéért) ---
// Most már nem kapnak 'tr' paramétert, és a DejaVu fontot használják

func renderHeaderPDF_UTF8(pdf *gofpdf.Fpdf, resume *data.ResumeData) {
	name, _ := resume.ContactInfo.Name.Get()
	email, _ := resume.ContactInfo.Email.Get()
	phone, _ := resume.ContactInfo.Phone.Get()
	website, _ := resume.ContactInfo.Website.Get()
	linkedin, _ := resume.ContactInfo.LinkedIn.Get()

	pdf.SetFont("DejaVu", "B", 18)                                      // Félkövér DejaVu
	pdf.CellFormat(contentWidthPDF, 10, name, "", 1, "C", false, 0, "") // Nincs tr()
	pdf.Ln(3)

	pdf.SetFont("DejaVu", "", 10) // Sima DejaVu
	contactParts := []string{}
	if email != "" {
		contactParts = append(contactParts, "Email: "+email)
	}
	if phone != "" {
		contactParts = append(contactParts, "Telefon: "+phone)
	}
	if website != "" {
		contactParts = append(contactParts, "Web: "+website)
	}
	if linkedin != "" {
		contactParts = append(contactParts, "LinkedIn: "+linkedin)
	}
	contactLine := strings.Join(contactParts, " | ")
	pdf.MultiCell(contentWidthPDF, lineHeightPDF, contactLine, "", "C", false) // Nincs tr()
	pdf.Ln(sectionSpacingPDF)
}

func renderSectionTitlePDF_UTF8(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("DejaVu", "B", 12) // Félkövér DejaVu
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(contentWidthPDF, lineHeightPDF*1.2, " "+title+" ", "", 1, "L", true, 0, "") // Nincs tr()
	pdf.Ln(lineHeightPDF / 2)
	pdf.SetFont("DejaVu", "", 10) // Vissza sima DejaVu-ra
}

func renderMultiLineTextPDF_UTF8(pdf *gofpdf.Fpdf, text string) {
	if text == "" {
		return
	}
	pdf.MultiCell(contentWidthPDF, lineHeightPDF, text, "", "L", false) // Nincs tr()
	pdf.Ln(lineHeightPDF / 2)
}

func renderExperienceEntryPDF_UTF8(pdf *gofpdf.Fpdf, entry *data.ExperienceEntry) {
	company, _ := entry.Company.Get()
	position, _ := entry.Position.Get()
	startDate, _ := entry.StartDate.Get()
	endDate, _ := entry.EndDate.Get()
	description, _ := entry.Description.Get()

	pdf.SetFont("DejaVu", "B", 10)                                                                                           // Félkövér DejaVu
	pdf.CellFormat(contentWidthPDF*0.7, lineHeightPDF, fmt.Sprintf("%s @ %s", position, company), "", 0, "L", false, 0, "")  // Nincs tr()
	pdf.SetFont("DejaVu", "", 10)                                                                                            // Sima DejaVu
	pdf.CellFormat(contentWidthPDF*0.3, lineHeightPDF, fmt.Sprintf("%s - %s", startDate, endDate), "", 1, "R", false, 0, "") // Nincs tr()

	renderMultiLineTextPDF_UTF8(pdf, description) // Ez már használja az új verziót
}

func renderEducationEntryPDF_UTF8(pdf *gofpdf.Fpdf, entry *data.EducationEntry) {
	institution, _ := entry.Institution.Get()
	degree, _ := entry.Degree.Get()
	startDate, _ := entry.StartDate.Get()
	endDate, _ := entry.EndDate.Get()
	description, _ := entry.Description.Get()

	pdf.SetFont("DejaVu", "B", 10)                                                                                            // Félkövér DejaVu
	pdf.CellFormat(contentWidthPDF*0.7, lineHeightPDF, fmt.Sprintf("%s @ %s", degree, institution), "", 0, "L", false, 0, "") // Nincs tr()
	pdf.SetFont("DejaVu", "", 10)                                                                                             // Sima DejaVu
	pdf.CellFormat(contentWidthPDF*0.3, lineHeightPDF, fmt.Sprintf("%s - %s", startDate, endDate), "", 1, "R", false, 0, "")  // Nincs tr()

	renderMultiLineTextPDF_UTF8(pdf, description) // Ez már használja az új verziót
}
