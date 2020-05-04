variable "test_string_list" {
  description = "list of strings"
  type = list(string)
  default = ["one", "two", "three"]
}

variable "test_number_list" {
  description = "list of numbers"
  type = list(number)
  default = [1, 2, 3]
}

variable "test_bool_list" {
  description = "list of bools"
  type = list(bool)
  default = [true, false]
}

variable "test_tuple_mv" {
  description = "multi-value tuple"
  type = tuple([string, number, bool])
  default = ["test", 1, true]
}

variable "test_string_map" {
  description = "test map for strings"
  type = map(string)
  default = {
    a = "ay"
    b = "bee"
    c = "cee"
  }
}

variable "test_object" {
  description = "test object"
  type = object({
    a=string,
    b=number,
    c=bool
  })
  default = {
    a = "ay",
    b = 10,
    c = false,
  }
}

variable "test_string_set" {
  description = "set of strings"
  type = set(string)
  default = ["one", "two", "three"]
}

variable "test_list_of_objects" {
  description = "test list of objects"
  type = list(object({
    a=string,
    b=number,
    c=bool
  }))
  default = [{
    a = "ay",
    b = 10,
    c = false,
  },
  {
    d = "dee",
    e = 20,
    f = true,
  }]
}

variable "test_object_with_list" {
  description = "test object with a list"
  type = object({
    a=list(string)
  })
  default = {
    a = ["a", "b", "c"]
  }
}