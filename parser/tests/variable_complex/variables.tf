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