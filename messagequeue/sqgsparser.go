package messagequeue

import (
	"errors"
	"reflect"

	"github.com/Pinablink/sqg/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

//
func prepareMap(mMap map[string]types.MessageAttributeValue, rType reflect.Type, key, value string) error {

	var iError error

	if rType.Name() == "string" {

		var iMessageAttributeValue types.MessageAttributeValue = types.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(value),
		}

		mMap[key] = iMessageAttributeValue

	} else {
		iError = errors.New(util.ERROR_TYPE_FIELD_STRUCT)
	}

	return iError
}

//
func ParserMessage(irefstruct interface{}) (map[string]types.MessageAttributeValue, error) {

	var iError error
	var mstruct reflect.Type = reflect.TypeOf(irefstruct)
	var nFields int = mstruct.NumField()
	var indexct int
	var mAttributes map[string]types.MessageAttributeValue

	if nFields > 0 {

		var mField reflect.StructField
		var mTag reflect.StructTag
		var mFieldType reflect.Type
		var mValue reflect.Value

		mAttributes = make(map[string]types.MessageAttributeValue)

		for indexct = 0; indexct < nFields; indexct++ {
			mField = mstruct.Field(indexct)
			mTag = mField.Tag
			mFieldType = mField.Type

			if len(mTag.Get("SQGS")) == 0 {
				iError = errors.New(util.NOT_FOUND_TAG_SQGS_IN_FIELD)
				break
			}

			mValue = reflect.ValueOf(irefstruct).FieldByName(mField.Name)
			iiError := prepareMap(mAttributes, mFieldType, mTag.Get("SQGS"), mValue.String())

			if iiError != nil {
				iError = iiError
				break
			}

		}

	} else {
		iError = errors.New(util.NOT_FOUND_FIELDS_CONTENT_STRUCT)
	}

	return mAttributes, iError
}
