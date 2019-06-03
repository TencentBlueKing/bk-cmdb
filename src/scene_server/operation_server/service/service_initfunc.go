package service

import "net/http"

func (s *Service) operation() {
	s.addAction(http.MethodPost, "/create/operation/chart", s.CreateStatisticChart, nil)
	s.addAction(http.MethodDelete, "/delete/operation/chart/{id}", s.DeleteStatisticChart, nil)
	s.addAction(http.MethodPost, "/update/operation/chart", s.UpdateStatisticChart, nil)
	s.addAction(http.MethodDelete, "/search/operation/chart", s.SearchStatisticCharts, nil)
	s.addAction(http.MethodPost, "/search/operation/chart/data", s.SearchChartData, nil)
}

func (s *Service) initService() {
	s.operation()
}
