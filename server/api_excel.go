package server

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/OnebookTechnology/jifengou/server/models"
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
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Length", fmt.Sprintf("%d", i.Size()))
	c := &BytesCounter{0}
	ServeContent(ctx.Writer, ctx.Request, fileName, f, c)
	return
}

// 平台券码导出
func ExportCoupon(ctx *gin.Context) {
	crossDomain(ctx)
	if err := CheckSession(ctx); err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid session. err:", err)
		return
	}

	pIdStr := ctx.Query("p_id")
	pId, err := strconv.Atoi(pIdStr)
	if err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid p_id. p_id:", pIdStr)
		return
	}
	p, err := server.DB.FindProductById(pId)
	if err != nil {
		if err == sql.ErrNoRows {
			sendSuccessResponse(ctx, nil)
			return
		}
		sendFailedResponse(ctx, Err, "FindAllExchangeRecord. err:", err)
		return
	}

	xlsx := excelize.NewFile()
	xlsx.SetCellValue("Sheet1", "A1", "商品ID")
	xlsx.SetCellValue("Sheet1", "B1", "商品名称")
	xlsx.SetCellValue("Sheet1", "A2", pId)
	xlsx.SetCellValue("Sheet1", "B2", p.ProductName)

	for status := models.CouponNotReleased; status <= models.CouponLogOut; status++ {
		records, err := server.DB.FindCouponsByProductId(pId, status, 1, 100000)
		if err != nil {
			if err == sql.ErrNoRows {
				sendSuccessResponse(ctx, nil)
				return
			}
			sendFailedResponse(ctx, Err, "FindAllExchangeRecord. err:", err)
			return
		}

		sheetName := models.CouponStatusMap[status+1]
		sheet := xlsx.NewSheet(sheetName)
		xlsx.SetCellValue(sheetName, "A1", "券码ID")
		xlsx.SetCellValue(sheetName, "B1", "平台券码")
		xlsx.SetCellValue(sheetName, "C1", "开始时间")
		xlsx.SetCellValue(sheetName, "D1", "结束时间")
		xlsx.SetCellValue(sheetName, "E1", "更新时间")
		xlsx.SetCellValue(sheetName, "F1", "商家券码")

		c := 0
		for _, record := range records {
			xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(c+2), record.CouponId)
			xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(c+2), record.CouponCode)
			xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(c+2), record.CouponStartTime)
			xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(c+2), record.CouponEndTime)
			xlsx.SetCellValue(sheetName, "E"+strconv.Itoa(c+2), record.UpdateTime)
			var bcodes string
			for _, s := range record.BCouponCodes {
				bcodes += s + "\n"
			}
			xlsx.SetCellValue(sheetName, "F"+strconv.Itoa(c+2), bcodes)
			c++
		}
		xlsx.SetActiveSheet(sheet)
	}

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
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Length", fmt.Sprintf("%d", i.Size()))
	c := &BytesCounter{0}
	ServeContent(ctx.Writer, ctx.Request, fileName, f, c)
	return
}

// 平台券码导出
func ExportBCoupon(ctx *gin.Context) {
	crossDomain(ctx)
	if err := CheckSession(ctx); err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid session. err:", err)
		return
	}

	pIdStr := ctx.Query("p_id")
	pId, err := strconv.Atoi(pIdStr)
	if err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid p_id. p_id:", pIdStr)
		return
	}
	xlsx := excelize.NewFile()
	p, err := server.DB.FindProductById(pId)
	if err != nil {
		if err == sql.ErrNoRows {
			sendSuccessResponse(ctx, nil)
			return
		}
		sendFailedResponse(ctx, Err, "FindAllExchangeRecord. err:", err)
		return
	}
	xlsx.SetCellValue("Sheet1", "A1", "商品ID")
	xlsx.SetCellValue("Sheet1", "B1", "商品名称")
	xlsx.SetCellValue("Sheet1", "A2", pId)
	xlsx.SetCellValue("Sheet1", "B2", p.ProductName)
	for status := models.CouponNotReleased; status <= models.CouponLogOut; status++ {
		records, err := server.DB.FindBCouponByStatus(status, pId, 1, 100000)
		if err != nil {
			if err == sql.ErrNoRows {
				sendSuccessResponse(ctx, nil)
				return
			}
			sendFailedResponse(ctx, Err, "FindAllExchangeRecord. err:", err)
			return
		}

		sheetName := models.CouponStatusMap[status+1]
		sheet := xlsx.NewSheet(sheetName)
		xlsx.SetCellValue(sheetName, "A1", "商家券码")
		xlsx.SetCellValue(sheetName, "B1", "开始时间")
		xlsx.SetCellValue(sheetName, "C1", "结束时间")
		xlsx.SetCellValue(sheetName, "D1", "更新时间")

		c := 0
		for _, record := range records {
			xlsx.SetCellValue(sheetName, "A"+strconv.Itoa(c+2), record.BCCode)
			xlsx.SetCellValue(sheetName, "B"+strconv.Itoa(c+2), record.BCStart)
			xlsx.SetCellValue(sheetName, "C"+strconv.Itoa(c+2), record.BCEnd)
			xlsx.SetCellValue(sheetName, "D"+strconv.Itoa(c+2), record.BCUpdateTime)
			c++
		}
		xlsx.SetActiveSheet(sheet)
	}

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
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	ctx.Header("Content-Length", fmt.Sprintf("%d", i.Size()))
	c := &BytesCounter{0}
	ServeContent(ctx.Writer, ctx.Request, fileName, f, c)
	return
}
