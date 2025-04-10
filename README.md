Feladat – Modern, szerkeszthető önéletrajz készítő alkalmazás

Cél:

Készíts egy modern, vizuálisan is letisztult önéletrajz-készítő alkalmazást Go vagy Python nyelven. A cél, hogy a felhasználó könnyen összeállíthassa az önéletrajzát, azt több formátumba is menthesse, illetve később szerkeszthesse is azt.

⸻

## Fő funkciók:

### Grafikus felület (GUI)
- Modern, átlátható design (pl. Material Design vagy hasonló)
- Drag & Drop mezők (pl. “Tapasztalat”, “Tanulmány”, “Készség” blokkok mozgatása)
- Élő előnézet (preview) mód
- Dark/Light mód váltás (opcionális)

### Adatbevitel
- Felhasználó megadhatja:
- Név, email, telefonszám, weboldal, LinkedIn
- Tanulmányok, munkatapasztalat, projektek
- Készségek, nyelvtudás, stb.
- Lehessen szakaszokat hozzáadni, eltávolítani, újrarendezni

### Mentés funkció (gombbal vagy menüből)
- Exportálás formátumok:
- PDF
- DOCX
- PNG (screenshot vagy render)
- Legyen egy saját fájlformátum (pl. .cvx) JSON vagy YAML alapú, amit később vissza lehet tölteni és szerkeszteni

### Betöltés és újraszerkesztés
- Felhasználó betölthet egy .cvx fájlt és ott folytathatja ahol abbahagyta

### Validálás
- Ne engedjen mentést, ha nincs megadva név vagy alapadat
- Legyen alapértelmezett sablon, ha semmit sem ír be

⸻

## Technikai elvárások:

### Nyelv:
- Go (pl. Fyne GUI framework)
vagy
- Python (pl. Tkinter, PyQt, vagy Kivy)

## Követelmények:
-  PDF export: pl. reportlab vagy wkhtmltopdf (Python), go-pdf (Go)
- DOCX export: pl. python-docx (Python), unioffice (Go)
- PNG export: képernyőkép vagy render a GUI-ból
- .cvx mentés: JSON alapú saját struktúra

## Egyéb:
- Kód legyen moduláris és tisztán szervezett
- Kommentált, jól dokumentált
- README legyen a repo-ban, amiben le van írva a build és a használat menete

⸻

## Extra (nem kötelező, de előny):
- Felhőalapú mentés (pl. Firebase vagy Google Drive API)
- Több nyelv támogatása
- Több CV sablon választható

⸻

## Leadási forma:
- GitHub repo
- Rövid videó/demo (pl. gif vagy képernyőfelvétel)
-  README.md (használati leírás, build lépések)’’’
