package main

// var (
// 	Result string
// )

// func TestRouteExtractParameter(t *testing.T) {
// 	// setup service
// 	ws := new(restful.WebService)
// 	ws.Consumes(restful.MIME_XML)
// 	ws.Route(ws.GET("/test/{param}").To(DummyHandler))
// 	restful.Add(ws)

// 	// setup request + writer
// 	bodyReader := strings.NewReader("<Sample><Value>42</Value></Sample>")
// 	httpRequest, _ := http.NewRequest("GET", "/test/THIS", bodyReader)
// 	httpRequest.Header.Set("Content-Type", restful.MIME_XML)
// 	httpWriter := httptest.NewRecorder()

// 	// run
// 	restful.DefaultContainer.ServeHTTP(httpWriter, httpRequest)

// 	if httpWriter.Code != 201 {
// 		fmt.Println(httpWriter.Code)
// 		t.Fatalf("WTF!")
// 	}
// 	if Result != "THIS" {
// 		t.Fatalf("Result is actually: %s", Result)
// 	}
// }

// func DummyHandler(rq *restful.Request, rp *restful.Response) {
// 	restful.DefaultResponseContentType(restful.MIME_JSON)

// 	Result = rq.PathParameter("param")
// 	rp.WriteHeader(123)
// 	user := User{"todo"}
// 	rp.WriteEntity(user)
// }
