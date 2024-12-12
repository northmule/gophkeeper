package data_type

const (
	// CardType данные банковских карт
	CardType = "card_type"
	// TextType произвольные текстовые данные
	TextType = "text_type"
	// BinaryType бинарные данные
	BinaryType = "binary_type"
	// MetaNameNote тип заметка для мета данных
	MetaNameNote = "meta_name_note"
	// MetaNameWebSite тип мета веб сайт
	MetaNameWebSite = "meta_name_website"
)

// TranslateDataType Тип поля в название
func TranslateDataType(dataType string) string {
	switch dataType {
	case CardType:
		return "Bank card details"
	case TextType:
		return "Text data"
	case BinaryType:
		return "Binary data"
	case MetaNameNote:
		return "Note"
	case MetaNameWebSite:
		return "Website"
	}

	return dataType
}
