package main

const articlesTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Список товаров</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f0f0f0;
            display: flex;
            flex-direction: column;
            align-items: center;
            margin: 0;
            padding: 20px 0;
        }
        div.article {
            width: 80%;
            margin-bottom: 20px;
            border: 1px solid #ddd;
            border-radius: 15px;
            background-color: #fff;
            box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
            word-wrap: break-word;
			padding: 5px 20px 20px 20px;
        }
		div.advert:hover{
			 background-color: #faf4ff;
		}
		p.descr {
			font-style: italic;
		}
        .image-slider {
            display: flex;
            overflow-x: auto;
			padding-bottom: 9px;
        }
        .image-slider img {
            max-height: 200px;
            margin-right: 10px;
        }
		.image-slider::-webkit-scrollbar {
		  width: 10px; /* Ширина полосы прокрутки */
		}
		
		.image-slider::-webkit-scrollbar-track {
		  background-color: #f1f1f1; /* Цвет фона полосы прокрутки */
		}
		
		.image-slider::-webkit-scrollbar-thumb {
		  background-color: #ccc; /* Цвет полосы прокрутки */
		}
		
		.image-slider::-webkit-scrollbar-thumb:hover {
		  background-color: #555; /* Цвет полосы прокрутки при наведении */
		}
    </style>
</head>
<body>
  {{range .}}
    <div class="article">
        <h3>{{ .Name}}</h3>
        <p class="descr">{{ .Brand}}</p>
        <p><b>Цена:</b> {{ .PriceSale}}</p>
        <p><b>Дата добавления: </b>{{.ParsingDate}}</p>
        <p><a href="{{.Link}}">{{.Link}}</a></p>
       	<div class="image-slider">
        {{range .Images}}
        	<img src="{{.}}" alt="Advert Image">
        {{end}}
</div>
    </div>
    {{end}}
</body>
</html>
`
