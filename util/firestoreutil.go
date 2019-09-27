package util

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

// MySnsData 外部から参照するSNS情報の構造体
type MySnsData struct {
	Facebook  string
	Twitter   string
	Instagram string
	Line      string
}

func getFirebaseClient(ctx context.Context) (*firestore.Client, error) {
	projectId := "snsqrcodegenerator-faad6"

	// mac 秘密鍵設置場所
	opt := option.WithCredentialsFile("/Users/snsqrcodegenerator-firebase-adminsdk-db.json")

	client, err := firestore.NewClient(ctx, projectId, opt)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// SaveUserItem Firestoreにデータを登録
func SaveUserItem(m MySnsData) (MySnsData, error) {
	ctx := context.Background()

	client, err := getFirebaseClient(ctx)
	if err != nil {
		return m, err
	}

	// "users" というコレクションの "1" というドキュメントに画面から渡されたデータを登録する
	_, err = client.Collection("users").Doc("1").Set(ctx, map[string]interface{}{
		"facebook":  m.Facebook,
		"twitter":   m.Twitter,
		"instagram": m.Instagram,
		"line":      m.Line,
	})

	if err != nil {
		return m, err
	}

	return m, nil
}

// Firestoreからデータを取得
func GetUserItem() (MySnsData, error) {
	ctx := context.Background()

	m := MySnsData{}

	client, err := getFirebaseClient(ctx)
	if err != nil {
		return m, err
	}

	doc := client.Collection("users").Doc("1")

	field, err := doc.Get(ctx)
	if err != nil {
		return m, err
	}

	data := field.Data()

	m.Facebook = data["facebook"].(string)
	m.Twitter = data["twitter"].(string)
	m.Instagram = data["instagram"].(string)
	m.Line = data["line"].(string)

	return m, nil
}

// Firestoreのデータを削除
func AllDiscard() error {
	ctx := context.Background()

	client, err := getFirebaseClient(ctx)
	if err != nil {
		return err
	}

	_, err = client.Collection("users").Doc("1").Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
