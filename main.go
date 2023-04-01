package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type Student struct {
	Name         string  `json:"name"`
	Age          int64   `json:"age"`
	AverageScore float64 `json:"average_score"`
}

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err
}

func main() {
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	//
	//newStudent := Student{
	//	Name:         "Gopher Test",
	//	Age:          10,
	//	AverageScore: 99.9,
	//}
	//
	//dataJSON, err := json.Marshal(newStudent)
	//js := string(dataJSON)
	//_, err = esclient.Index().
	//	Index("students").
	//	BodyJson(js).
	//	Do(ctx)
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println("[Elastic][InsertProduct]Insertion Successful")

	var students []Student

	searchSource := elastic.NewSearchSource()
	//searchSource.Query(elastic.NewMatchQuery("name", "Goph"))
	searchSource.Query(elastic.NewWildcardQuery("name", "*op*"))

	searchService := esclient.Search().Index("students").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		var student Student
		err := json.Unmarshal(hit.Source, &student)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		students = append(students, student)
	}

	if err != nil {
		fmt.Println("Fetching student fail: ", err)
	} else {
		for _, s := range students {
			fmt.Printf("Student found Name: %s, Age: %d, Score: %f \n", s.Name, s.Age, s.AverageScore)
		}
	}
}
