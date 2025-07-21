package main

import (
	"log"
	"net/http"
	"os"

	"yoyaku/auth"
	"yoyaku/db"
	"yoyaku/handler"
	"yoyaku/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

func main() {

	//データベースとの接続
	sqlDB, err := utils.NewDBConnection()
	if err != nil {
		log.Fatalf("データベースに接続できませんでした: %v", err)
	}
	defer sqlDB.Close()
	// sql.DBからsqlcのクエリオブジェクトを生成
	queries := db.New(sqlDB)

	// 1. OAuth設定の初期化
	if err := auth.Setup(); err != nil {
		log.Fatalf("OAuth設定の初期化に失敗しました: %v", err)
	}

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("環境変数 SECRET_KEY が設定されていません")
	}

	// セッション情報を保存するためのストア (キーは秘密の値にしてください)
	var store = sessions.NewCookieStore([]byte(secretKey))

	// Ginのルーティング
	r := gin.Default()

	// CORS設定
	// フロントエンドのURL(http://localhost:3000)からのリクエストを許可
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"}
	config.AllowCredentials = true // Cookieの送受信を許可
	r.Use(cors.New(config))

	// セッションミドルウェア
	// 全てのリクエストでセッションストアを利用可能にする
	r.Use(func(c *gin.Context) {
		c.Set("session_store", store)
		c.Next()
	})

	// 3. ルーティングの設定
	r.GET("/login", handler.HandleGoogleLogin)
	r.GET("/callback", func(c *gin.Context) {
		handler.HandleGoogleCallback(c, queries)
	})

	// フロントエンドがユーザー情報を確認するためのAPIエンドポイント
	// api := r.Group("/api")
	// {
	// 	api.GET("/me", handler.HandleGetMe)
	// 	api.POST("/logout", handler.HandleLogout)
	// 	api.POST("/reservations", func(c *gin.Context) { handler.Handlereservations(c, queries) })
	// 	api.GET("/reservations/me", func(c *gin.Context) { handler.HandlereservationsMe(c, queries) })
	// 	api.PUT("/reservations/cancel", func(c *gin.Context) { handler.HandlereservationsCancele(c, queries) })
	// }
	api := r.Group("/api")
	{
		// ユーザー認証関連
		api.GET("/me", handler.HandleGetMe)
		api.POST("/logout", handler.HandleLogout)

		// 予約関連のAPIをグループ化
		reservations := api.Group("/reservations")
		{
			// POST /api/reservations
			// 新しい予約を作成
			reservations.POST("", func(c *gin.Context) {
				handler.Handlereservations(c, queries)
			})

			reservations.PUT("", func(c *gin.Context) {
				handler.HandlereservationsEdit(c, queries)
			})

			// GET /api/reservations/me
			// ログインユーザー自身の予約一覧を取得
			reservations.GET("/me", func(c *gin.Context) {
				handler.HandlereservationsMe(c, queries)
			})

			// PUT /api/reservations/cancel
			// 予約をキャンセル
			reservations.PUT("/cancel", func(c *gin.Context) {
				handler.HandlereservationsCancele(c, queries)
			})

			// GET /api/reservations?month=... や ?date=...
			// クエリパラメータに応じて全ユーザーの予約を期間で絞り込んで取得
			reservations.GET("", func(c *gin.Context) {
				if c.Query("month") != "" {
					handler.HandlerListByMonth(c, queries)
				} else if c.Query("start") != "" && c.Query("end") != "" {
					handler.HandlerListByWeek(c, queries)
				} else if c.Query("date") != "" {
					handler.HandlerListByDate(c, queries)
				} else {
					// ここに全件取得や、パラメータがない場合のエラー処理などを記述
					c.JSON(http.StatusBadRequest, gin.H{"error": "有効なクエリパラメータがありません"})
				}
			})
		}
	}

	// 4. サーバーの起動
	log.Println("Started server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
