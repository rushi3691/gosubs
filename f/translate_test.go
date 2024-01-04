package f

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/text/language"
)

func TestTranslate(t *testing.T) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api_key := os.Getenv("API_KEY")

	translator, err := NewTranslator("en", api_key, language.English, language.Spanish)
	if err != nil {
		log.Fatalf("Failed to create translator: %v", err)
	}

	translated, err := translator.Translate("Hello, world!")
	log.Println(translated)
	if err != nil {
		log.Fatalf("Failed to translate text: %v", err)
	}

	// fmt.Println(translated)
	if translated != "¡Hola Mundo!" {
		t.Errorf("Expected %q, got %q", "¡Hola Mundo!", translated)
	}

	t.Log("TestTranslate passed")
}
