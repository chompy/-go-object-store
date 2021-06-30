package main

//type QueryOperator string

type Query struct {
	Operator int
	Not      bool
	Tags     []string
}

// "tag-id-(>1000)"
// "tag-time-(<12/31/2001)"
// "tag-test-* AND tag-test-1 OR "

/*func parseQuery(queryString string) {

}
*/
