package main

import (
	"log"
	"net/http"
	"yoyaku/auth"
	"yoyaku/handler"
	"yoyaku/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func main() {
	type Data1 struct {
		gorm.Model
		Title   string
		Content string
	}

	//データベースとの接続
	db, err := utils.NewDBConnection()
	if err != nil {
		println("データベースの接続できませんでした")
	}

	// 1. OAuth設定の初期化
	if err := auth.Setup(); err != nil {
		log.Fatalf("OAuth設定の初期化に失敗しました: %v", err)
	}

	// セッション情報を保存するためのストア (キーは秘密の値にしてください)
	var store = sessions.NewCookieStore([]byte("your-super-secret-key"))

	err = db.AutoMigrate(&Data1{})
	if err != nil {
		println("Userのマイグレーションに失敗しました。")
	}

	// Ginのルーティング
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		db.Create(&Data1{Title: "Test", Content: "test"})
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello Database",
		})
	})

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
	r.GET("/callback", handler.HandleGoogleCallback)

	// フロントエンドがユーザー情報を確認するためのAPIエンドポイント
	api := r.Group("/api")
	{
		api.GET("/me", handler.HandleGetMe)
		api.POST("/logout", handler.HandleLogout)
	}

	// 4. サーバーの起動
	log.Println("Started server on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
