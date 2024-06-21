package lang

var enUS = &Language{
	name: "en-US",
	lines: map[string]string{
		"malformed-request":            "Malformed request",
		"malformed-json":               "Malformed JSON",
		"auth.invalid-credentials":     "Invalid credentials.",
		"auth.no-credentials-provided": "Invalid or missing authentication header.",
		"auth.jwt-invalid":             "Your authentication token is invalid.",
		"auth.jwt-not-valid-yet":       "Your authentication token is not valid yet.",
		"auth.jwt-expired":             "Your authentication token is expired.",
	},
	validation: validationLines{
		rules: map[string]string{
			"required":                           "The :field is required.",
			"required.element":                   "The :field elements are required.",
			"float32":                            "The :field must be numeric.",
			"float32.element":                    "The :field elements must be numeric.",
			"float64":                            "The :field must be numeric.",
			"float64.element":                    "The :field elements must be numeric.",
			"int":                                "The :field must be an integer.",
			"int.element":                        "The :field elements must be integers.",
			"int8":                               "The :field must be an integer.",
			"int8.element":                       "The :field elements must be integers.",
			"int16":                              "The :field must be an integer.",
			"int16.element":                      "The :field elements must be integers.",
			"int32":                              "The :field must be an integer.",
			"int32.element":                      "The :field elements must be integers.",
			"int64":                              "The :field must be an integer.",
			"int64.element":                      "The :field elements must be integers.",
			"uint":                               "The :field must be a positive integer.",
			"uint.element":                       "The :field elements must be positive integers.",
			"uint8":                              "The :field must be a positive integer.",
			"uint8.element":                      "The :field elements must be positive integers.",
			"uint16":                             "The :field must be a positive integer.",
			"uint16.element":                     "The :field elements must be positive integers.",
			"uint32":                             "The :field must be a positive integer.",
			"uint32.element":                     "The :field elements must be positive integers.",
			"uint64":                             "The :field must be a positive integer.",
			"uint64.element":                     "The :field elements must be positive integers.",
			"string":                             "The :field must be a string.",
			"string.element":                     "The :field elements must be strings.",
			"array":                              "The :field must be an array.",
			"array.element":                      "The :field elements must be arrays.",
			"min.string":                         "The :field must be at least :min characters.",
			"min.numeric":                        "The :field must be at least :min.",
			"min.array":                          "The :field must have at least :min items.",
			"min.file":                           "The :field must be at least :min KiB.",
			"min.object":                         "The :field must have at least :min fields.",
			"min.string.element":                 "The :field elements must be at least :min characters.",
			"min.numeric.element":                "The :field elements must be at least :min.",
			"min.array.element":                  "The :field elements must have at least :min items.",
			"min.object.element":                 "The :field elements must have at least :min fields.",
			"max.string":                         "The :field may not have more than :max characters.",
			"max.numeric":                        "The :field may not be greater than :max.",
			"max.array":                          "The :field may not have more than :max items.",
			"max.file":                           "The :field may not be greater than :max KiB.",
			"max.object":                         "The :field may not have more than :max fields.",
			"max.string.element":                 "The :field elements may not have more than :max characters.",
			"max.numeric.element":                "The :field elements may not be greater than :max.",
			"max.array.element":                  "The :field elements may not have more than :max items.",
			"max.object.element":                 "The :field elements may not have more than :max fields.",
			"between.string":                     "The :field must be between :min and :max characters.",
			"between.numeric":                    "The :field must be between :min and :max.",
			"between.array":                      "The :field must have between :min and :max items.",
			"between.object":                     "The :field must have between :min and :max fields.",
			"between.file":                       "The :field must be between :min and :max KiB.",
			"between.string.element":             "The :field elements must be between :min and :max characters.",
			"between.numeric.element":            "The :field elements must be between :min and :max.",
			"between.array.element":              "The :field elements must have between :min and :max items.",
			"between.object.element":             "The :field elements must have between :min and :max fields.",
			"greater_than.string":                "The :field must be longer than the :other.",
			"greater_than.numeric":               "The :field must be greater than the :other.",
			"greater_than.array":                 "The :field must have more items than the :other.",
			"greater_than.file":                  "The :field must be larger than the :other.",
			"greater_than.object":                "The :field must have more fields than the :other.",
			"greater_than.string.element":        "The :field elements must be longer than the :other.",
			"greater_than.numeric.element":       "The :field elements must be greater than the :other.",
			"greater_than.array.element":         "The :field elements must have more items than the :other.",
			"greater_than.object.element":        "The :field elements must have more fields than the :other.",
			"greater_than_equal.string":          "The :field must be longer or have the same length as the :other.",
			"greater_than_equal.numeric":         "The :field must be greater or equal to the :other.",
			"greater_than_equal.array":           "The :field must have more or the same amount of items as the :other.",
			"greater_than_equal.file":            "The :field must be the same size or larger than the :other.",
			"greater_than_equal.object":          "The :field must have at least as many fields as the :other.",
			"greater_than_equal.string.element":  "The :field elements must be longer or have the same length as the :other.",
			"greater_than_equal.numeric.element": "The :field elements must be greater or equal to the :other.",
			"greater_than_equal.array.element":   "The :field elements must have more or the same amount of items as the :other.",
			"greater_than_equal.object.element":  "The :field elements must have at least as many fields as the :other.",
			"lower_than.string":                  "The :field must be shorter than the :other.",
			"lower_than.numeric":                 "The :field must be lower than the :other.",
			"lower_than.array":                   "The :field must have less items than the :other.",
			"lower_than.file":                    "The :field must be smaller than the :other.",
			"lower_than.object":                  "The :field must have less fields than the :other.",
			"lower_than.string.element":          "The :field elements must be shorter than the :other.",
			"lower_than.numeric.element":         "The :field elements must be lower than the :other.",
			"lower_than.array.element":           "The :field elements must have less items than the :other.",
			"lower_than.object.element":          "The :field elements must have less fields than the :other.",
			"lower_than_equal.string":            "The :field must be shorter or have the same length as the :other.",
			"lower_than_equal.numeric":           "The :field must be lower or equal to the :other.",
			"lower_than_equal.array":             "The :field must have less or the same amount of items as the :other.",
			"lower_than_equal.file":              "The :field must be the same size or smaller than the :other.",
			"lower_than_equal.object":            "The :field must have at most as many fields as the :other.",
			"lower_than_equal.string.element":    "The :field elements must be shorter or have the same length as the :other.",
			"lower_than_equal.numeric.element":   "The :field elements must be lower or equal to the :other.",
			"lower_than_equal.array.element":     "The :field elements must have less or the same amount of items as the :other.",
			"lower_than_equal.object.element":    "The :field elements must have at most as many fields as the :other.",
			"distinct":                           "The :field must have only distinct values.",
			"distinct.element":                   "The :field elements must have only distinct values.",
			"digits":                             "The :field must be digits only.",
			"digits.element":                     "The :field elements must be digits only.",
			"regex":                              "The :field format is invalid.",
			"regex.element":                      "The :field element format is invalid.",
			"email":                              "The :field must be a valid email address.",
			"email.element":                      "The :field elements must be valid email addresses.",
			"size.string":                        "The :field must be exactly :value characters-long.",
			"size.numeric":                       "The :field must be exactly :value.",
			"size.array":                         "The :field must contain exactly :value items.",
			"size.file":                          "The :field must be exactly :value KiB.",
			"size.object":                        "The :field must have exactly :value fields.",
			"size.string.element":                "The :field elements must be exactly :value characters-long.",
			"size.numeric.element":               "The :field elements must be exactly :value.",
			"size.array.element":                 "The :field elements must contain exactly :value items.",
			"size.object.element":                "The :field elements must have exactly :value fields.",
			"alpha":                              "The :field may only contain letters.",
			"alpha.element":                      "The :field elements may only contain letters.",
			"alpha_dash":                         "The :field may only contain letters, numbers, dashes and underscores.",
			"alpha_dash.element":                 "The :field elements may only contain letters, numbers, dashes and underscores.",
			"alpha_num":                          "The :field may only contain letters and numbers.",
			"alpha_num.element":                  "The :field elements may only contain letters and numbers.",
			"starts_with":                        "The :field must start with one of the following values: :values.",
			"starts_with.element":                "The :field elements must start with one of the following values: :values.",
			"ends_with":                          "The :field must end with one of the following values: :values.",
			"ends_with.element":                  "The :field elements must end with one of the following values: :values.",
			"doesnt_start_with":                  "The :field must not start with any of the following values: :values.",
			"doesnt_start_with.element":          "The :field elements must not start with any of the following values: :values.",
			"in":                                 "The :field must have one of the following values: :values.",
			"in.element":                         "The :field elements must have one of the following values: :values.",
			"not_in":                             "The :field must not have one of the following values: :values.",
			"not_in.element":                     "The :field elements must not have one of the following values: :values.",
			"in_field":                           "The :field must exist in the :other.",
			"in_field.element":                   "The :field elements must exist in the :other.",
			"not_in_field":                       "The :field must not exist in the :other.",
			"not_in_field.element":               "The :field elements must not exist in the :other.",
			"timezone":                           "The :field must be a valid time zone.",
			"timezone.element":                   "The :field elements must be valid time zones.",
			"ip":                                 "The :field must be a valid IP address.",
			"ip.element":                         "The :field elements must be valid IP addresses.",
			"ipv4":                               "The :field must be a valid IPv4 address.",
			"ipv4.element":                       "The :field elements must be valid IPv4 addresses.",
			"ipv6":                               "The :field must be a valid IPv6 address.",
			"ipv6.element":                       "The :field elements must be valid IPv6 addresses.",
			"json":                               "The :field must be a valid JSON string.",
			"json.element":                       "The :field elements must be valid JSON strings.",
			"url":                                "The :field must be a valid URL.",
			"url.element":                        "The :field elements must be valid URLs.",
			"uuid":                               "The :field must be a valid UUID.",
			"uuid.element":                       "The :field elements must be valid UUIDs.",
			"bool":                               "The :field must be a boolean.",
			"bool.element":                       "The :field elements must be booleans.",
			"same":                               "The :field and the :other must match.",
			"same.element":                       "The :field elements and the :other must match.",
			"different":                          "The :field and the :other must be different.",
			"different.element":                  "The :field elements and the :other must be different.",
			"file":                               "The :field must be a file.",
			"mime":                               "The :field must be a file of type: :values.",
			"image":                              "The :field must be an image.",
			"extension":                          "The :field must be a file with one of the following extensions: :values.",
			"file_count":                         "The :field must have exactly :value file(s).",
			"min_file_count":                     "The :field must have at least :value file(s).",
			"max_file_count":                     "The :field may not have more than :value file(s).",
			"file_count_between":                 "The :field must have between :min and :max files.",
			"date":                               "The :field is not a valid date.",
			"date.element":                       "The :field elements are not valid dates.",
			"before":                             "The :field must be a date before :date.",
			"before.element":                     "The :field elements must be dates before :date.",
			"before_equal":                       "The :field must be a date before or equal to :date.",
			"before_equal.element":               "The :field must be dates before or equal to :date.",
			"after":                              "The :field must be a date after :date.",
			"after.element":                      "The :field elements must be dates after :date.",
			"after_equal":                        "The :field must be a date after or equal to :date.",
			"after_equal.element":                "The :field elements must be dates after or equal to :date.",
			"date_equals":                        "The :field must be a date equal to :date.",
			"date_equals.element":                "The :field elements must be dates equal to :date.",
			"object":                             "The :field must be an object.",
			"object.element":                     "The :field elements must be objects.",
			"unique":                             "The :field has already been taken.",
			"unique.element":                     "The :field element value has already been taken.",
			"exists":                             "The :field does not exist.",
			"exists.element":                     "The :field element value does not exist.",
			"keysin":                             "The :field keys must be one of the following: :values.",
			"keysin.element":                     "The :field elements keys must be one of the following: :values.",
		},
		fields: map[string]string{
			"":        "body",
			"email":   "email address",
			"perPage": "number of records per page",
			"*":       "object property",
		},
	},
}

// SetDefaultLine set the language line identified by the given key in the
// default "en-US" language.
// Values set this way can be overridden by language files.
func SetDefaultLine(key, line string) {
	enUS.lines[key] = line
}

// SetDefaultValidationRule set the validation error message for the rule identified by
// the given key in the default "en-US" language.
// Values set this way can be overridden by language files.
func SetDefaultValidationRule(key, line string) {
	enUS.validation.rules[key] = line
}

// SetDefaultFieldName set the field name used in validation error message placeholders
// for the given field in the default "en-US" language.
// Values set this way can be overridden by language files.
func SetDefaultFieldName(field, name string) {
	enUS.validation.fields[field] = name
}
