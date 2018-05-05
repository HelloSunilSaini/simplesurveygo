from jsonschema import validate
import json

schema0 = {
            "type" : "object",
            "properties" : {
                "colors" : {"type" : [{"type" : "object",
                                        "properties" : {
                                            "color" : {"type": "string"},
                                            "catogory" : {"type": "string"},
                                            "type" : {"type" : "string"},
                                            "code" : {"type" : "object",
                                                    "properties" : {
                                                        "rgba" : {"type" : ["int"]},
                                                        "hex" : {"type" : "string"}
                                                    }
                                            }
                                        }
                }]}
            }
}

a = {
        "colors": [
        {
            "color": "black",
            "category": "hue",
            "type": "primary",
            "code": {
            "rgba": [255,255,255,1],
            "hex": "#000"
            }
        },
        {
            "color": "white",
            "category": "value",
            "code": {
            "rgba": [0,0,0,1],
            "hex": "#FFF"
            }
        },
        {
            "color": "red",
            "category": "hue",
            "type": "primary",
            "code": {
            "rgba": [255,0,0,1],
            "hex": "#FF0"
            }
        }
        ]
    }

def validate_schema(payload):
    print (payload)
    try:
        validate(payload,schema0)
        print (True)
    except:
        print (False)


b = json.dumps(a)
validate_schema(b)
validate_schema(a)

