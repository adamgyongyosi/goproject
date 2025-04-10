// io/storage.go
package io

import (
	"encoding/json"
	"fmt"
	"os"
	"resume-builder/data" // Győződj meg róla, hogy a modul neve helyes!
)

// Mentési struktúra a kapcsolathoz
type saveContactInfo struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
}

// Mentési struktúra a tapasztalathoz
type saveExperienceEntry struct {
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Description string `json:"description,omitempty"`
}

// Mentési struktúra a tanulmányokhoz
type saveEducationEntry struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
	Description string `json:"description,omitempty"`
}

// Mentési struktúra a teljes önéletrajzhoz
type saveResumeData struct {
	ContactInfo saveContactInfo       `json:"contactInfo"`
	Summary     string                `json:"summary,omitempty"`
	Experience  []saveExperienceEntry `json:"experience,omitempty"`
	Education   []saveEducationEntry  `json:"education,omitempty"` // ÚJ MEZŐ
	Skills      []string              `json:"skills,omitempty"`    // ÚJ MEZŐ (string slice)

}

// SaveResume menti az önéletrajz adatait JSON formátumban.
func SaveResume(resume *data.ResumeData, filePath string) error {
	// Contact adatok másolása
	contactData := saveContactInfo{}
	contactData.Name, _ = resume.ContactInfo.Name.Get()
	contactData.Email, _ = resume.ContactInfo.Email.Get()
	contactData.Phone, _ = resume.ContactInfo.Phone.Get()
	contactData.Website, _ = resume.ContactInfo.Website.Get()
	contactData.LinkedIn, _ = resume.ContactInfo.LinkedIn.Get()

	// Fő mentési struktúra létrehozása
	saveData := saveResumeData{
		ContactInfo: contactData,
		Summary:     "", // Kezdőérték
	}
	// Summary másolása
	summaryVal, _ := resume.Summary.Get()
	saveData.Summary = summaryVal

	// Experience lista másolása
	experienceItems := []saveExperienceEntry{}
	expListInterface, errExpGet := resume.Experience.Get()
	if errExpGet != nil {
		fmt.Println("Hiba az Experience lista lekérésekor mentésnél:", errExpGet)
	} else {
		for _, itemInterface := range expListInterface {
			if entry, ok := itemInterface.(*data.ExperienceEntry); ok {
				saveEntry := saveExperienceEntry{}
				saveEntry.Company, _ = entry.Company.Get()
				saveEntry.Position, _ = entry.Position.Get()
				saveEntry.StartDate, _ = entry.StartDate.Get()
				saveEntry.EndDate, _ = entry.EndDate.Get()
				saveEntry.Description, _ = entry.Description.Get()
				experienceItems = append(experienceItems, saveEntry)
			} else {
				fmt.Printf("Hiba mentésnél: Az Experience listában nem várt típus található: %T\n", itemInterface)
			}
		}
	}
	saveData.Experience = experienceItems

	// Education lista másolása
	educationItems := []saveEducationEntry{}
	eduListInterface, errEduGet := resume.Education.Get()
	if errEduGet != nil {
		fmt.Println("Hiba az Education lista lekérésekor mentésnél:", errEduGet)
	} else {
		for _, itemInterface := range eduListInterface {
			if entry, ok := itemInterface.(*data.EducationEntry); ok {
				saveEntry := saveEducationEntry{}
				saveEntry.Institution, _ = entry.Institution.Get()
				saveEntry.Degree, _ = entry.Degree.Get()
				saveEntry.StartDate, _ = entry.StartDate.Get()
				saveEntry.EndDate, _ = entry.EndDate.Get()
				saveEntry.Description, _ = entry.Description.Get()
				educationItems = append(educationItems, saveEntry)
			} else {
				fmt.Printf("Hiba mentésnél: Education listában nem várt típus: %T\n", itemInterface)
			}
		}
	}
	saveData.Education = educationItems

	// ÚJ: Skills lista mentése
	skillsVal, errSkillsGet := resume.Skills.Get()
	if errSkillsGet != nil {
		fmt.Println("Hiba a Skills lista lekérésekor mentésnél:", errSkillsGet)
		saveData.Skills = []string{} // Üres lista mentése hiba esetén
	} else {
		saveData.Skills = skillsVal // String slice közvetlen hozzárendelése
	}

	// JSON generálás és fájlba írás
	jsonData, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// LoadResume betölti az önéletrajz adatait JSON fájlból.
func LoadResume(filePath string) (*data.ResumeData, error) {
	// Fájl olvasása
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// JSON értelmezése a segédstruktúrába
	var saveData saveResumeData
	err = json.Unmarshal(jsonData, &saveData)
	if err != nil {
		return nil, err
	}

	// Új, köthető ResumeData példány létrehozása
	newResume := data.NewResumeData()

	// ContactInfo és Summary adatok beállítása
	newResume.ContactInfo.Name.Set(saveData.ContactInfo.Name)
	newResume.ContactInfo.Email.Set(saveData.ContactInfo.Email)
	newResume.ContactInfo.Phone.Set(saveData.ContactInfo.Phone)
	newResume.ContactInfo.Website.Set(saveData.ContactInfo.Website)
	newResume.ContactInfo.LinkedIn.Set(saveData.ContactInfo.LinkedIn)
	newResume.Summary.Set(saveData.Summary)

	// Experience lista betöltése és beállítása
	loadedExperienceEntries := []interface{}{}
	for _, saveEntry := range saveData.Experience {
		newEntry := data.NewExperienceEntry()
		newEntry.Company.Set(saveEntry.Company)
		newEntry.Position.Set(saveEntry.Position)
		newEntry.StartDate.Set(saveEntry.StartDate)
		newEntry.EndDate.Set(saveEntry.EndDate)
		newEntry.Description.Set(saveEntry.Description)
		loadedExperienceEntries = append(loadedExperienceEntries, newEntry)
	}
	errExpSet := newResume.Experience.Set(loadedExperienceEntries)
	if errExpSet != nil {
		fmt.Printf("Hiba az Experience lista beállításakor betöltésnél: %v\n", errExpSet)
	}

	// Education lista betöltése és beállítása
	loadedEducationEntries := []interface{}{}
	for _, saveEntry := range saveData.Education {
		newEntry := data.NewEducationEntry()
		newEntry.Institution.Set(saveEntry.Institution)
		newEntry.Degree.Set(saveEntry.Degree)
		newEntry.StartDate.Set(saveEntry.StartDate)
		newEntry.EndDate.Set(saveEntry.EndDate)
		newEntry.Description.Set(saveEntry.Description)
		loadedEducationEntries = append(loadedEducationEntries, newEntry)
	}
	errEduSet := newResume.Education.Set(loadedEducationEntries)
	if errEduSet != nil {
		fmt.Printf("Hiba az Education lista beállításakor betöltésnél: %v\n", errEduSet)
	}

	// ÚJ: Skills lista betöltése
	// A saveData.Skills már egy []string, ezt közvetlenül be tudjuk állítani
	errSkillsSet := newResume.Skills.Set(saveData.Skills)
	if errSkillsSet != nil {
		fmt.Printf("Hiba a Skills lista beállításakor betöltésnél: %v\n", errSkillsSet)
	}

	return newResume, nil
}
