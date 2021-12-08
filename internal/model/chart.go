package model

type PaginateData struct {
	Links struct {
		Pagination struct {
			Total       int `json:"total,omitempty"`
			PerPage     int `json:"per_page,omitempty"`
			CurrentPage int `json:"current_page,omitempty"`
			From        int `json:"from,omitempty"`
			To          int `json:"to,omitempty"`
		} `json:"pagination"`
	} `json:"links"`

	Data []interface{} `json:"data"`
}

type LineChartData struct {
	Labels   []string  `json:"labels,omitempty"`
	Datasets []Dataset `json:"dataset,omitempty"`
}

type BarChartData struct {
	Labels   []string  `json:"labels,omitempty"`
	Datasets []Dataset `json:"dataset,omitempty"`
}

type PieChartData struct {
	Labels   []string  `json:"labels,omitempty"`
	Datasets []Dataset `json:"dataset,omitempty"`
}

type Dataset struct {
	Label string `json:"label,omitempty"`
	Name  string `json:"name,omitempty"`
	Data  []int  `json:"data,omitempty"`
}
