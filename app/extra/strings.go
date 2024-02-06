package strings

func AllLower(word string) bool {
	
	for ( _,r := range word ) {
   		if unicode.IsUpper(r) {
      			return false
  		 }
	}

	return true
}
