package WildBerries

type WbAdvertsJson struct {
	State   int `json:"state"`
	Version int `json:"version"`
	Params  struct {
		Version        int    `json:"version"`
		Curr           string `json:"curr"`
		Spp            int    `json:"spp"`
		PayloadVersion int    `json:"payloadVersion"`
	} `json:"params"`
	Data struct {
		Products []struct {
			Sort            int     `json:"__sort"`
			Ksort           int     `json:"ksort"`
			Time1           int     `json:"time1"`
			Time2           int     `json:"time2"`
			Dist            int     `json:"dist"`
			Id              int     `json:"id"`
			Root            int     `json:"root"`
			KindId          int     `json:"kindId"`
			SubjectId       int     `json:"subjectId"`
			SubjectParentId int     `json:"subjectParentId"`
			Name            string  `json:"name"`
			Brand           string  `json:"brand"`
			BrandId         int     `json:"brandId"`
			SiteBrandId     int     `json:"siteBrandId"`
			Supplier        string  `json:"supplier"`
			SupplierId      int     `json:"supplierId"`
			Sale            int     `json:"sale"`
			PriceU          int     `json:"priceU"`
			SalePriceU      int     `json:"salePriceU"`
			LogisticsCost   int     `json:"logisticsCost"`
			SaleConditions  int     `json:"saleConditions"`
			ReturnCost      int     `json:"returnCost"`
			Pics            int     `json:"pics"`
			Rating          int     `json:"rating"`
			ReviewRating    float64 `json:"reviewRating"`
			Feedbacks       int     `json:"feedbacks"`
			PanelPromoId    int     `json:"panelPromoId,omitempty"`
			PromoTextCard   string  `json:"promoTextCard,omitempty"`
			PromoTextCat    string  `json:"promoTextCat,omitempty"`
			Volume          int     `json:"volume"`
			ViewFlags       int     `json:"viewFlags"`
			Colors          []struct {
				Name string `json:"name"`
				Id   int    `json:"id"`
			} `json:"colors"`
			Sizes []struct {
				Name       string `json:"name"`
				OrigName   string `json:"origName"`
				Rank       int    `json:"rank"`
				OptionId   int    `json:"optionId"`
				ReturnCost int    `json:"returnCost"`
				Wh         int    `json:"wh"`
				Sign       string `json:"sign"`
				Payload    string `json:"payload"`
			} `json:"sizes"`
			DiffPrice bool `json:"diffPrice"`
			Log       struct {
				Cpm           int `json:"cpm,omitempty"`
				Promotion     int `json:"promotion,omitempty"`
				PromoPosition int `json:"promoPosition,omitempty"`
				Position      int `json:"position,omitempty"`
				AdvertId      int `json:"advertId,omitempty"`
			} `json:"log"`
		} `json:"products"`
	} `json:"data"`
}

// https://card.wb.ru/cards/v1/detail?appType=1&curr=rub&spp=30&nm=25968980;75563253

type FullAdvertJson struct {
	State  int `json:"state"`
	Params struct {
		Version        int    `json:"version"`
		Curr           string `json:"curr"`
		Spp            int    `json:"spp"`
		PayloadVersion int    `json:"payloadVersion"`
	} `json:"params"`
	Data struct {
		Products []struct {
			Id              int    `json:"id"`
			Root            int    `json:"root"`
			KindId          int    `json:"kindId"`
			SubjectId       int    `json:"subjectId"`
			SubjectParentId int    `json:"subjectParentId"`
			Name            string `json:"name"`
			Brand           string `json:"brand"`
			BrandId         int    `json:"brandId"`
			SiteBrandId     int    `json:"siteBrandId"`
			Supplier        string `json:"supplier"`
			SupplierId      int    `json:"supplierId"`
			PriceU          int    `json:"priceU"`
			SalePriceU      int    `json:"salePriceU"`
			LogisticsCost   int    `json:"logisticsCost"`
			Sale            int    `json:"sale"`
			Extended        struct {
				BasicSale    int `json:"basicSale"`
				BasicPriceU  int `json:"basicPriceU"`
				ClientSale   int `json:"clientSale"`
				ClientPriceU int `json:"clientPriceU"`
			} `json:"extended"`
			SaleConditions int     `json:"saleConditions"`
			ReturnCost     int     `json:"returnCost"`
			Pics           int     `json:"pics"`
			Rating         int     `json:"rating"`
			ReviewRating   float64 `json:"reviewRating"`
			Feedbacks      int     `json:"feedbacks"`
			PanelPromoId   int     `json:"panelPromoId"`
			PromoTextCard  string  `json:"promoTextCard"`
			PromoTextCat   string  `json:"promoTextCat"`
			Volume         int     `json:"volume"`
			ViewFlags      int     `json:"viewFlags"`
			Colors         []struct {
				Name string `json:"name"`
				Id   int    `json:"id"`
			} `json:"colors"`
			Promotions []int `json:"promotions"`
			Sizes      []struct {
				Name       string        `json:"name"`
				OrigName   string        `json:"origName"`
				Rank       int           `json:"rank"`
				OptionId   int           `json:"optionId"`
				ReturnCost int           `json:"returnCost"`
				Stocks     []interface{} `json:"stocks"`
			} `json:"sizes"`
			DiffPrice bool `json:"diffPrice"`
		} `json:"products"`
	} `json:"data"`
}
