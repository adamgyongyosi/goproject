// data/resume.go
package data

import "fyne.io/fyne/v2/data/binding"

// Kapcsolattartási információk
type ContactInfo struct {
    Name     binding.String `json:"name" yaml:"name"`
    Email    binding.String `json:"email" yaml:"email"`
    Phone    binding.String `json:"phone" yaml:"phone"`
    Website  binding.String `json:"website,omitempty" yaml:"website,omitempty"`
    LinkedIn binding.String `json:"linkedin,omitempty" yaml:"linkedin,omitempty"`
}

// Munkatapasztalat bejegyzés struktúrája
type ExperienceEntry struct {
    Company     binding.String `json:"company" yaml:"company"`
    Position    binding.String `json:"position" yaml:"position"`
    StartDate   binding.String `json:"startDate" yaml:"startDate"`
    EndDate     binding.String `json:"endDate" yaml:"endDate"`
    Description binding.String `json:"description,omitempty" yaml:"description,omitempty"`
}

// Segédfüggvény új ExperienceEntry létrehozásához
func NewExperienceEntry() *ExperienceEntry {
    entry := &ExperienceEntry{
        Company:     binding.NewString(),
        Position:    binding.NewString(),
        StartDate:   binding.NewString(),
        EndDate:     binding.NewString(),
        Description: binding.NewString(),
    }
    entry.Company.Set("")
    entry.Position.Set("")
    entry.StartDate.Set("")
    entry.EndDate.Set("")
    entry.Description.Set("")
    return entry
}

// Tanulmány bejegyzés struktúrája
type EducationEntry struct {
    Institution binding.String `json:"institution" yaml:"institution"`
    Degree      binding.String `json:"degree" yaml:"degree"`
    StartDate   binding.String `json:"startDate" yaml:"startDate"`
    EndDate     binding.String `json:"endDate" yaml:"endDate"`
    Description binding.String `json:"description,omitempty" yaml:"description,omitempty"`
}

// Segédfüggvény új EducationEntry létrehozásához
func NewEducationEntry() *EducationEntry {
    entry := &EducationEntry{
        Institution: binding.NewString(),
        Degree:      binding.NewString(),
        StartDate:   binding.NewString(),
        EndDate:     binding.NewString(),
        Description: binding.NewString(),
    }
    entry.Institution.Set("")
    entry.Degree.Set("")
    entry.StartDate.Set("")
    entry.EndDate.Set("")
    entry.Description.Set("")
    return entry
}

// A teljes önéletrajz adatszerkezete
type ResumeData struct {
    ContactInfo ContactInfo         `json:"contactInfo" yaml:"contactInfo"`
    Summary     binding.String      `json:"summary,omitempty" yaml:"summary,omitempty"`
    Experience  binding.UntypedList `json:"-" yaml:"-"`
    Education   binding.UntypedList `json:"-" yaml:"-"` // Education lista hozzáadva
    Skills      binding.StringList  `json:"-" yaml:"-"` // ÚJ LISTA (StringList)
}

// Funkció egy új, üres önéletrajz adatpéldány létrehozásához
// Figyelj, hogy a 'ContactInfo: contact,' sor benne legyen!
func NewResumeData() *ResumeData {
    // ContactInfo inicializálása
    contact := ContactInfo{
        Name:     binding.NewString(),
        Email:    binding.NewString(),
        Phone:    binding.NewString(),
        Website:  binding.NewString(),
        LinkedIn: binding.NewString(),
    }
    contact.Name.Set("")
    contact.Email.Set("")
    contact.Phone.Set("")
    contact.Website.Set("")
    contact.LinkedIn.Set("")

    // Summary inicializálása
    summary := binding.NewString()
    summary.Set("")

    // Experience lista inicializálása
    experienceList := binding.NewUntypedList()

    // Education lista inicializálása
    educationList := binding.NewUntypedList()

    skillsList := binding.NewStringList()

    // A VISSZAADOTT STRUKTÚRA (Ellenőrizd, hogy ez a rész pontosan így van-e!)
    return &ResumeData{
        ContactInfo: contact,
        Summary:     summary,
        Experience:  experienceList,
        Education:   educationList,
        Skills:      skillsList, // ÚJ MEZŐ hozzáadása
    }
}