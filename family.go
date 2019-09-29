package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type person struct {
	name   string //  person.name never used can be removed
	gender string
}

type parent struct {
	male   string
	female string
}

func NewPerson(Name string, Gender string) person {
	p := person{name: Name, gender: Gender}
	return p
}

func NewParents(dad string, mom string) parent {
	p := parent{male: dad, female: mom}
	return p
}

var persons = make(map[string]person)  // key : name of person,  value: person
var parents = make(map[string]parent)  // key : name of person , value: parent of that person {dad,mom}
var childs = make(map[string][]string) // key : name of person,  value: slice of childs
var couples = make(map[string]string)  // key : name of first person, value: name of second person
var rootMother = "Queen-Anga"

func Message(msg string, silent bool) { // silent argument is for when first time build family don't print anyting
	if !silent {
		fmt.Println(msg)
	}
}

func NotExists(present bool, exists bool, silent bool) { // to check whether value of query (GET_RELATIONSHIP) is empty or not

	if !present && exists {
		Message("NONE", silent)
	} else if !exists {
		Message("PERSON_NOT_FOUND", silent)
	} else {
		Message("", silent)
	}
}

/*

AddChild() create a person then append that person into list of childs of mother

*/

func AddChild(mom string, name string, gender string, silent bool) {
	_, exists := persons[name]
	if exists {
		prnts, exists := parents[name]
		if exists && prnts.female == mom && prnts.male == couples[mom] {
			Message("CHILD_ADDITION_FAILED", silent)
			return
		}
	}
	m, exists := persons[mom]
	if exists && m.gender == "Male" { //  mom can't be Male
		Message("CHILD_ADDITION_FAILED", silent)
		return
	}
	if !exists && mom != rootMother {
		Message("PERSON_NOT_FOUND", silent)
		return
	}
	prsn := NewPerson(name, gender)        // create new person
	prnts := NewParents(couples[mom], mom) // create parents for that person
	parents[name] = prnts
	childs[mom] = append(childs[mom], name) // append this person into list of mom's child
	persons[name] = prsn
	Message("CHILD_ADDITION_SUCCEEDED", silent)
}

/*

AddCouple() creates a person if one of p1,p2 is not in family tree
and assign partner of p1 as p2 and vice versa

*/

func AddCouple(p1 string, p2 string) {

	_, p1_exists := persons[p1]
	_, p2_exists := persons[p2]
	if !p1_exists && !p2_exists {
		return
	} else if p1_exists && !p2_exists {
		gender := persons[p1].gender
		if gender == "Male" {
			gender = "Female"
		} else {
			gender = "Male"
		}
		prsn := NewPerson(p2, gender)
		persons[p2] = prsn
	} else if p2_exists && !p1_exists {
		gender := persons[p2].gender
		if gender == "Male" {
			gender = "Female"
		} else {
			gender = "Male"
		}
		prsn := NewPerson(p1, gender)
		persons[p1] = prsn
	}
	couples[p1] = p2
	couples[p2] = p1

}

/*
   To get Son or Daughter of a person
   GetChild() iterates list of child of a person if she is female otherwise iterates in his wife's child
*/

func GetChild(prsn string, gender string, silent bool) {

	p, exists := persons[prsn]
	spouse, married := couples[prsn]
	present := false

	if exists && !married && p.gender == "Male" { // if person is male and doesn't have wife then he would have no child
		NotExists(present, exists, silent)
		return
	}

	if exists {
		key := prsn
		if p.gender == "Male" {
			key = spouse
		}
		for _, child := range childs[key] {
			if persons[child].gender == gender {
				fmt.Print(child, " ")
				present = true
			}
		}

	}
	NotExists(present, exists, silent)
}

/*

GetSibling() iterate in mom's childs

*/

func GetSibling(prsn string, silent bool) {
	_, exists := persons[prsn]
	prnts, exist := parents[prsn]
	present := false
	if !exist {
		NotExists(present, exists, silent)
		return
	}
	for _, sibling := range childs[prnts.female] { // iterate in mom's childs
		if sibling != prsn {
			fmt.Print(sibling, " ")
			present = true
		}
	}
	NotExists(present, exists, silent)
}

/*
 To get paternal Aunt and Uncle

 GetPaternal() iterates in list of child of dad's mother

*/

func GetPaternal(prsn string, gender string, silent bool) {
	_, exists := persons[prsn]
	prnts, exist := parents[prsn]
	present := false
	if !exist {
		NotExists(present, exists, silent)
		return
	}
	grand_parent, exist := parents[prnts.male]
	if !exist {
		NotExists(present, exists, silent)
		return
	}

	for _, child := range childs[grand_parent.female] { // iterate in grand_parent's childs
		if (child != prnts.male && child != prnts.female) && persons[child].gender == gender {
			fmt.Print(child, " ")
			present = true
		}
	}
	NotExists(present, exists, silent)

}

/*
 To get Maternal Aunt and Uncle
 GetMaternal() iterates in list of child of mom's mother

*/
func GetMaternal(prsn string, gender string, silent bool) {
	_, exists := persons[prsn]
	prnts, exist := parents[prsn]
	present := false
	if !exist {
		NotExists(present, exists, silent)
		return
	}
	grand_parent, exist := parents[prnts.female]
	if !exist {
		NotExists(present, exists, silent)
		return
	}

	for _, child := range childs[grand_parent.female] {
		if (child != prnts.male && child != prnts.female) && persons[child].gender == gender {
			fmt.Print(child, " ")
			present = true
		}
	}
	NotExists(present, exists, silent)

}

/*
 To get  Brother-In-Law and Sister-In-Law
 GetinLaw() iterates in spouse's sibling if person is married and also  iterates in own siblings if sibling is married

*/
func GetinLaw(prsn string, gender string, silent bool) {
	_, exists := persons[prsn]
	present := false

	if !exists {
		NotExists(present, exists, silent)
		return
	}
	spouse, married := couples[prsn]

	if married {
		prnts, exists := parents[spouse]
		if exists {
			for _, sibling := range childs[prnts.female] {
				if sibling != spouse && persons[sibling].gender == gender {
					fmt.Print(sibling, " ")
					present = true
				}
			}
		}
	}

	prnts, exist := parents[prsn]

	if !exist {
		NotExists(present, exists, silent)
		return
	}
	for _, sibling := range childs[prnts.female] {
		spouse, married := couples[sibling]
		if married && persons[spouse].gender == gender {
			present = true
			fmt.Print(spouse, " ")
		}
	}
	NotExists(present, exists, silent)
}

/*
 To Process commands from file

*/

func commands(path string, silent bool) {

	fileHandle, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer fileHandle.Close()
	fileScanner := bufio.NewScanner(fileHandle)

	for fileScanner.Scan() {
		str := fileScanner.Text()
		cmd := strings.Fields(str)
		if len(cmd) < 3 {
			continue
		}
		switch cmd[0] {
		case "ADD_CHILD":
			AddChild(cmd[1], cmd[2], cmd[3], silent)
			break
		case "GET_RELATIONSHIP":
			if cmd[2] == "Maternal-Aunt" {
				GetMaternal(cmd[1], "Female", silent)
			} else if cmd[2] == "Maternal-Uncle" {
				GetMaternal(cmd[1], "Male", silent)
			} else if cmd[2] == "Paternal-Uncle" {
				GetPaternal(cmd[1], "Male", silent)
			} else if cmd[2] == "Maternal-Aunt" {
				GetPaternal(cmd[1], "Female", silent)
			} else if cmd[2] == "Brother-In-Law" {
				GetinLaw(cmd[1], "Male", silent)
			} else if cmd[2] == "Sister-In-Law" {
				GetinLaw(cmd[1], "Female", silent)
			} else if cmd[2] == "Siblings" {
				GetSibling(cmd[1], silent)
			} else if cmd[2] == "Son" {
				GetChild(cmd[1], "Male", silent)
			} else if cmd[2] == "Daughter" {
				GetChild(cmd[1], "Female", silent)
			}
			break
		case "ADD_COUPLE":
			AddCouple(cmd[1], cmd[2])
		}
	}
}
