# Modifiers

Given the path contains

Key | Value
--------
a|x
b|y
c|z
d|j

## valuesText

This gets implicitly invoked, but can be explicitly called.

~a~b~c -> xyz
~a~b~c|valuesText ->xyz

## base64

Base64 encodes the text

## jsonlist
Just takes the values and makes a list
~a~b~c|jsonlist -> ['x', 'y', 'z']

## jsonkeyedobject
Uses the keys you specfied into vault as the keys to their values in the object
~a~b~c|jsonkeyedobject -> {'a':'x', 'b':'y', 'c':'z'}

## jsonpairedobject
Takes pairs of values and makes the first into a key, second into a value. An odd number of keys is an error.
~a~b~c~d|jsonpairedobject -> {'x':'y', 'z':'j'}

## jsonobject2list(name,value)
Takes a json object and splits it into a list of objects with the keys coming from the call, and the values from the object.
~a~b~c|jsonkeyedobject|jsonobject2list(name, value)[{'name':'a', 'value':'x'}, {'name':'b', 'value':'y'}, {'name':'c', 'value':'z'}]
~a~b|jsonpairedobject|jsonobject2list(user, password)[{'user':'x', 'password','y'}]
~a~b~c~d|jsonpairedobject|jsonobject2list(user, password)[{'user':'x', 'password','y'}, {'user':'z', 'password','j'}]

## json2htaccess
Takes an object list (jsonobject2list) and uses the key called user and the key called password from each object to make an htaccess file. Any object in the list without these keys is an error.


