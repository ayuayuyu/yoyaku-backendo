package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"yoyaku/db"
	"yoyaku/types"
	"yoyaku/utils"

	"github.com/gin-gonic/gin"
)

func Handlereservations(c *gin.Context, queries *db.Queries) {
	var req types.ReservationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "リクエストの形式が正しくありません",
		})
		return
	}
	println("タイトル", req.Title, "開始時間", req.StartTime.String(), "終了時間", req.EndTime.String()) //デバック用

	// ユーティリティ関数を呼び出す
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		return
	}

	// 重複チェック
	count, err := queries.CheckOverlappingReservation(context.Background(), db.CheckOverlappingReservationParams{
		StartTime: req.EndTime,
		EndTime:   req.StartTime,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "重複チェック中にエラーが発生しました"})
		return
	}
	fmt.Printf("count : %v", count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "この時間帯には既に予約があります"})
		return
	}

	// 登録処理
	_, err = queries.CreateReservation(context.Background(), db.CreateReservationParams{
		UserID:    userID,
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "予約の登録に失敗しました"})
		return
	}
	reservation, err := queries.GetReservationLastInserted(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約情報の取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"id":         reservation.ID,
		"user_id":    reservation.UserID,
		"title":      reservation.Title,
		"start_time": reservation.StartTime,
		"end_time":   reservation.EndTime,
		"created_at": reservation.CreatedAt,
		"updated_at": reservation.UpdatedAt,
	})
}

func HandlereservationsMe(c *gin.Context, queries *db.Queries) {
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		return
	}

	reservations, err := queries.ListReservationsByUserID(context.Background(), userID)
	if err != nil {
		log.Println("予約取得エラー:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約の取得に失敗しました"})
		return
	}

	println("user_id", userID)
	println("reservations", reservations)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reservations,
	})
}

func HandlereservationsCancele(c *gin.Context, queries *db.Queries) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが指定されていません"})
		return
	}

	// 文字列のIDをuint64に変換
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 変換に失敗した場合（例：IDが数値ではない）
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDの形式が正しくありません"})
		return
	}
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		return
	}

	err = queries.CanceledReservationByID(context.Background(), db.CanceledReservationByIDParams{
		UserID: userID,
		ID:     id,
	})
	if err != nil {
		log.Println("予約キャンセルエラー:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約のキャンセルに失敗しました"})
		return
	}

	println("user_id", userID)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Reservation Canceled",
	})
}

func HandlereservationsEdit(c *gin.Context, queries *db.Queries) {
	idStr := c.Query("id")
	var req types.ReservationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "リクエストの形式が正しくありません",
		})
		return
	}
	println("タイトル", req.Title, "開始時間", req.StartTime.String(), "終了時間", req.EndTime.String()) //デバック用
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDが指定されていません"})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDの形式が正しくありません"})
		return
	}
	userID, ok := utils.GetUserIDFromSession(c)
	if !ok {
		return
	}

	err = queries.UpdateReservationByID(context.Background(), db.UpdateReservationByIDParams{
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		ID:        id,
	})
	if err != nil {
		log.Println("予約編集エラー:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約の編集に失敗しました"})
		return
	}

	println("user_id", userID)
	updated, err := queries.GetReservationByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新後の予約取得に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"id":         updated.ID,
		"user_id":    updated.UserID,
		"title":      updated.Title,
		"start_time": updated.StartTime,
		"end_time":   updated.EndTime,
		"created_at": updated.CreatedAt,
		"updated_at": updated.UpdatedAt,
	})
}

func HandlerListByMonth(c *gin.Context, queries *db.Queries) {
	monthStr := c.Query("month") // "2025-07" のような文字列
	if monthStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monthクエリパラメータは必須です"})
		return
	}

	// "YYYY-MM" 形式の文字列をtime.Timeオブジェクトに変換
	layout := "2006-01"
	t, err := time.Parse(layout, monthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "monthの形式が正しくありません (YYYY-MM)"})
		return
	}

	// 月の初日と最終日の翌日を計算
	startOfMonth := t
	endOfMonth := t.AddDate(0, 1, 0)

	reservations, err := queries.ListReservationsByMonth(context.Background(), db.ListReservationsByMonthParams{
		EndTime:   startOfMonth,
		StartTime: endOfMonth,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約の取得に失敗しました"})
		return
	}

	// 予約が存在しない場合にnullではなく空配列を返す
	if reservations == nil {
		reservations = []db.ListReservationsByMonthRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reservations,
	})
}

func HandlerListByWeek(c *gin.Context, queries *db.Queries) {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "startとendクエリパラメータは必須です"})
		return
	}

	layout := "2006-01-02"
	startTime, err1 := time.Parse(layout, startStr)
	endTime, err2 := time.Parse(layout, endStr)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日付の形式が正しくありません (YYYY-MM-DD)"})
		return
	}

	endTime = endTime.AddDate(0, 0, 1)

	reservations, err := queries.ListReservationsByWeek(context.Background(), db.ListReservationsByWeekParams{
		Starttime: startTime,
		Endtime:   endTime,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約の取得に失敗しました"})
		return
	}

	if reservations == nil {
		reservations = []db.ListReservationsByWeekRow{}
	}

	fmt.Printf("reservations : %v", reservations)

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reservations,
	})
}

func HandlerListByDate(c *gin.Context, queries *db.Queries) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "dateクエリパラメータは必須です"})
		return
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "日付の形式が正しくありません (YYYY-MM-DD)"})
		return
	}

	startOfDay := date
	endOfDay := date.AddDate(0, 0, 1)

	reservations, err := queries.ListReservationsByDate(context.Background(), db.ListReservationsByDateParams{
		StartTime: endOfDay,
		EndTime:   startOfDay,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "予約の取得に失敗しました"})
		return
	}

	if reservations == nil {
		reservations = []db.ListReservationsByDateRow{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   reservations,
	})
}
