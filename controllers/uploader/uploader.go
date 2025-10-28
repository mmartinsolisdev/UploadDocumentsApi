package uploader

import (
	"UploadDocumentsAPI/database"
	"UploadDocumentsAPI/models"
	_ "bytes"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	_ "log"
	_ "mime/multipart"
	_ "os"
	 "path/filepath"
	"strconv"
	_ "strings"
	_ "net/http"

	"github.com/gofiber/fiber/v2"
)

func GetMembershipsList(c *fiber.Ctx) error {

	Id := c.Query("Id")
	DocType := c.Query("DocType")
	Language := c.Query("Language")
	DocName := c.Query("DocName")
	SaleType := c.Query("SaleType")
	ContractCode := c.Query("ContractCode")
	code, err := strconv.Atoi(ContractCode)

	if err != nil {
		fmt.Println("La conversión no se puedo realizar")
	}
	//fmt.Println(DocType)

	//Database connection
	db := database.DBConn
	//contractTexts := &models.ContractTexts{}
	//Se Ejecuta la consulta y se almacena en la variable results
	data := models.ContractTexts{
		Id:           Id,
		Language:     Language,
		DocType:      DocType,
		DocName:      DocName,
		SaleType:      SaleType,
		ContractCode: code,
	}
	var results []models.ContractTexts
	//db.Table("ContractTexts").Find(&results)
	//db.Model(contractTexts).Where("cxla = ?", Language).Find(&results)
	db.Where(data).Select("cxID", "cxDocType", "cxla", "cxDocName", "cxSaleType","cxContractCode").Find(&results)
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
	db.Model(contractTexts).Select("cxSaleType").Distinct().Pluck("cxSaleType", &combosData)
	combos["SaleType"] = combosData
	db.Model(salesRoom).Select("srContractCode, srID").Find(&salesRoom)
	combos["Code"] = salesRoom

	return c.JSON(combos)
}

func UploadFile(c *fiber.Ctx) error {

	Id := c.Query("Id")
	DocType := c.Query("DocType")
	Language := c.Query("Language")
	DocName := c.Query("DocName")
	SaleType := c.Query("SaleType")
	ContractCode := c.Query("ContractCode")
	code, err := strconv.Atoi(ContractCode)

	if err != nil {
		fmt.Println("La conversión no se puedo realizar")
		return err
	}

	db := database.DBConn
	data := models.ContractTexts{
		Id:           Id,
		Language:     Language,
		DocType:      DocType,
		DocName:      DocName,
		SaleType:      SaleType,
		ContractCode: code,
	}
	contractTexts := &models.ContractTexts{}

	//Acepta el archivo como multipart form dentro del parametro documents
	file, err := c.FormFile("documents")
	if err != nil {
		fmt.Println("The file could not be obtained.")
		return c.SendStatus(404)
	}

	//Se abre el archivo y se almacena en memoria
	openedfile, err := file.Open()

	if err != nil {
		fmt.Println("Error to read the file")
		return err
	}
	//Se cierra el archivo al final de la función
	defer openedfile.Close()
	// fmt.Println(openfile)
	//Se lee el archivo con ioutil.ReadAll y se almacena el contenido en Bytes
	fileBytes, err := ioutil.ReadAll(openedfile)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var fileExtension = filepath.Ext(file.Filename)
	//fmt.Println(extension)
	if fileExtension != ".docx" {
		return c.SendStatus(415) //415 Unsupported Media Type
	}
	// Update que adjunta el archivo en bytes al campo cxTextBinary
	db.Model(contractTexts).Where(data).Update("cxTextBinary", fileBytes)
	return c.SendStatus(200) //c.SendStatus(200) //c.JSON("contents")
}
