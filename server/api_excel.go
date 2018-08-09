package server

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"os"
	"strconv"
)

func ExportUser(ctx *gin.Context) {
	crossDomain(ctx)
	if err := CheckSession(ctx); err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid session. err:", err)
		return
	}
	records, err := server.DB.FindAllExchangeRecord(1, 100000)
	if err != nil {
		if err == sql.ErrNoRows {
			sendSuccessResponse(ctx, nil)
			return
		}
		sendFailedResponse(ctx, Err, "FindAllExchangeRecord. err:", err)
		return
	}
	xlsx := excelize.NewFile()
	sheet := xlsx.NewSheet("Sheet1")
	xlsx.SetCellValue("Sheet1", "A1", "手机号")
	xlsx.SetCellValue("Sheet1", "B1", "兑换商品ID")
	xlsx.SetCellValue("Sheet1", "C1", "兑换商品")
	xlsx.SetCellValue("Sheet1", "D1", "兑换时间")
	xlsx.SetCellValue("Sheet1", "E1", "平台券")
	xlsx.SetCellValue("Sheet1", "F1", "消费券")

	for i, record := range records {
		xlsx.SetCellValue("Sheet1", "A"+strconv.Itoa(i+2), record.PhoneNumber)
		xlsx.SetCellValue("Sheet1", "B"+strconv.Itoa(i+2), record.PId)
		xlsx.SetCellValue("Sheet1", "C"+strconv.Itoa(i+2), record.Name)
		xlsx.SetCellValue("Sheet1", "D"+strconv.Itoa(i+2), record.ExTime)
		xlsx.SetCellValue("Sheet1", "E"+strconv.Itoa(i+2), record.PCode)
		xlsx.SetCellValue("Sheet1", "F"+strconv.Itoa(i+2), record.BCodes)
	}
	xlsx.SetActiveSheet(sheet)
	fileName := "ex_" + nowTimestampString() + ".xlsx"
	filePath := "/root/online/jifengou/images/" + fileName
	xlsx.SaveAs(filePath)
	// Save xlsx file by the given path.
	f, err := os.Open(filePath)
	if err != nil {
		sendFailedResponse(ctx, Err, "file load err:", err)
		return
	}
	defer f.Close()
	i, err := f.Stat()
	if err != nil {
		sendFailedResponse(ctx, Err, "file Stat err:", err)
		return
	}
	ctx.Header("Content-Type", "application/x-download")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Length", fmt.Sprintf("%d", i.Size()))
	c := &BytesCounter{0}
	ServeContent(ctx.Writer, ctx.Request, fileName, f, c)
	return
}
