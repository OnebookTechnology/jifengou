package server

import (
	"database/sql"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
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
	// Save xlsx file by the given path.
	ctx.Header("Content-Type", "application/vnd.ms-excel")
	xlsx.Write(ctx.Writer)
}
