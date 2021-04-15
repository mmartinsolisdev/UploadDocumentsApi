package uploader

import (
	 "fmt"
	 "io/ioutil"
 	 _"bytes"
	 _"log"
	 "UploadDocumentsAPI/database"
	 "UploadDocumentsAPI/models"
		_"encoding/json"
		_"strings"
		"strconv"
		_"os"
		_"mime/multipart"

	"github.com/gofiber/fiber/v2"
)


func GetMembershipsList(c *fiber.Ctx) error {
	//db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
// SELECT * FROM users WHERE name = "jinzhu";

	Id := c.Query("Id")
	DocType := c.Query("DocType")
	Language := c.Query("Language")
	DocName := c.Query("DocName")
	ContractCode := c.Query("ContractCode")
	code, err := strconv.Atoi(ContractCode)

	if err != nil {
		fmt.Println("La conversión no se puedo realizar")
	}

				fmt.Println(string(Id))
				fmt.Println(DocType)
				fmt.Println(Language)
				fmt.Println(DocName)
				//fmt.Println(ContractCodee)

	//Database connection
	db := database.DBConn
	//contractTexts := &models.ContractTexts{}
	//Se Ejecuta la consulta y se almacena en la variable results
	data := models.ContractTexts {
		Id: Id,
		Language: Language,
		DocType: DocType,
		DocName: DocName,
		ContractCode: code,
	}
	var results []models.ContractTexts
	//db.Table("ContractTexts").Find(&results)
	//db.Model(contractTexts).Where("cxla = ?", Language).Find(&results)
	db.Where(data).Select("cxID", "cxDocType", "cxla", "cxDocName", "cxContractCode").Find(&results)
	return c.JSON(results)
}



func GetCombos(c *fiber.Ctx) error {
	//Database connection
	db := database.DBConn
	contractTexts := &models.ContractTexts{}
	salesRoom := []models.SalesRooms{}
	//Se Ejecuta la consulta y se almacena en la variable results

	//var codigos = []Codigos{}
	//var results map[string]interface{}
	var combosData []string
	var combos = make(map[string]interface{})
	//"cxID", "cxDocType", "cxla", "cxDocName", "cxContractCode"
	db.Model(contractTexts).Select("cxID").Distinct().Pluck("cxId", &combosData)
	combos["Ids"] = combosData
	db.Model(contractTexts).Select("cxDocName").Distinct().Pluck("cxDocName", &combosData)
	combos["DocName"] = combosData
	db.Model(contractTexts).Select("cxla").Distinct().Pluck("cxLa", &combosData)
	combos["Language"] = combosData
	db.Model(contractTexts).Select("cxDocType").Distinct().Pluck("cxDocType", &combosData)
	combos["DocType"] = combosData
	db.Model(salesRoom).Select("srContractCode, srID").Find(&salesRoom)
	combos["Code"] =  salesRoom
	//fmt.Println(salesRoom)


/*	data := [][]string{
    {"Perro", "Pez", "Gato"},
    {"Nissan", "Ford", "Honda"},
}
data = append(data, []string{"GoLang", "C#", "PHP"})*/


/*cities := []Codigos{
	Codigos{Code: "1",Text: "Newport"},
	Codigos{Code: "2",Text: "Vistacay"},
}*/
//cities = append(cities, item2)
//fmt.Println(cities)
//combos["Code"] =  cities
/*scores := map[string]string{"Code": "1", "Text": "Newport"}
	fmt.Println(scores)*/
	return c.JSON(combos)
}





func UploadFile(c *fiber.Ctx) error {

	Id := c.Query("Id")
	DocType := c.Query("DocType")
	Language := c.Query("Language")
	DocName := c.Query("DocName")
	ContractCode := c.Query("ContractCode")
	code, err := strconv.Atoi(ContractCode)

	if err != nil {
		fmt.Println("La conversión no se puedo realizar")
		return err
	}

	db := database.DBConn
	data := models.ContractTexts {
		Id: Id,
		Language: Language,
		DocType: DocType,
		DocName: DocName,
		ContractCode: code,
	}
		contractTexts := &models.ContractTexts{}

/*
		// Get first file from form field "document":
		if form, err := c.MultipartForm(); err == nil {
    // Get all files from "documents" key:
    files := form.File["documents"]
    // Loop through files:
    for _, file := range files {
			//fmt.Println(file)
			//fmt.Println(file)
			// Save the files to disk:
			openfile, err := file.Open()
			fileData, err := ioutil.ReadAll(openfile)

			fmt.Println(fileData)
			filerc, err := os.Open(fileData)

	 	 if err != nil{
	 		log.Fatal(err)
	 		}
	 	defer filerc.Close()
	 		buf := new(bytes.Buffer)
	    buf.ReadFrom(filerc)
	    contents := buf.Bytes()
			//fmt.Println(contents)
			db.Model(contractTexts).Where(data).Update("cxTextBinary", contents)
    }
  }
*/
	/*
  //file, err := c.FormFile("document")
	//fmt.Print(file)
	 filerc, err := os.Open(file.Filename)
	 if err != nil{
		log.Fatal(err)
	}
	defer filerc.Close()
		buf := new(bytes.Buffer)
   buf.ReadFrom(filerc)
   contents := buf.Bytes()

   fmt.Print(contents)
	 // Update with conditions
	db.Model(contractTexts).Where(data).Update("cxTextBinary", contents)
*/


//Acepta el archivo como multipart form dentro del parametro documents
//if file, err := c.FormFile("documents"); err == nil {


	file, err := c.FormFile("documents")
	 if err != nil {
		 fmt.Println("No se pudo obtener el archivo")
		 return c.SendStatus(404)
	 }

//Se abre el archivo y se almacena en memoria
	 openedfile, err := file.Open()
//Se cierra el archivo al final de la función
	 defer openedfile.Close()
	// fmt.Println(openfile)
	 if err != nil {
		 fmt.Println("Error al abrir el archivo")
		 return err
	 }

	 //Se lee el archivo con ioutil.ReadAll y se almacena el contenido en Bytes
	 fileBytes, err := ioutil.ReadAll(openedfile)
	 if err != nil {
    fmt.Println(err)
		return err
		}

// Update que adjunta el archivo en bytes al campo cxTextBinary
db.Model(contractTexts).Where(data).Update("cxTextBinary", fileBytes)
//}

   // Save file to root directory:
   //return c.SaveFile(file, fmt.Sprintf("./%s", file.Filename))
	// Parse the multipart form:
   /*if form, err := c.MultipartForm(); err == nil {
     // Get all files from "documents" key:
     files := form.File["documents"]
		 //out, _ := os.Open(header.Filename)
     // => []*multipart.FileHeader
     // Loop through files:


     for _, file := range files {
       //fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			 //fmt.Println(file)


	//fmt.Printf("%s", b)
		 db.Model(contractTexts).Where(data).Update("cxTextBinary", file)
       // Save the files to disk:
       if err := c.SaveFile(file, fmt.Sprintf("./%s", file.Filename)); err != nil {
         return err
       }
     }
     return err
   }*/
	 return c.SendStatus(200) //c.JSON("contents")
}
