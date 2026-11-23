package domain

type PublicAnalytics struct {
	Overview   PublicAnalyticsOverview   `json:"overview"`
	Status     []PublicAnalyticsStatus   `json:"status"`
	Aging      []PublicAnalyticsAging    `json:"aging"`
	Throughput []PublicAnalyticsDay      `json:"throughput"`
	DueSoon    []PublicAnalyticsDeadline `json:"due_soon"`
}

type PublicAnalyticsOverview struct {
	Total              int64    `json:"Всего"`
	Completed          int64    `json:"Выполнено"`
	Open               int64    `json:"Открыто"`
	Overdue            int64    `json:"Просрочено"`
	DueSoon            int64    `json:"Срок в ближайшие 7 дней"`
	CompletionRate     *float64 `json:"Доля выполнения (%)"`
	AverageCompletionH *float64 `json:"Среднее время выполнения (ч)"`
}

type PublicAnalyticsStatus struct {
	Status string `json:"Статус"`
	Count  int64  `json:"Задачи"`
}

type PublicAnalyticsAging struct {
	Bucket string `json:"Просрочка"`
	Count  int64  `json:"Задачи"`
}

type PublicAnalyticsDay struct {
	Date      string `json:"День"`
	Created   int64  `json:"Создано"`
	Completed int64  `json:"Закрыто"`
}

type PublicAnalyticsDeadline struct {
	Date  string `json:"День"`
	Count int64  `json:"Задачи"`
}
