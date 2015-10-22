Research to see if it's possible to extend a SS namespace definition to the point
where the DO APIs can be consumed directly without a plugin.

## Computing a response resource collection type

In the case of namespace `get` actions CWF does not know about the resource type a-priori. It needs
to infer the type from the response payload. The algorithm is as follows:

1. `kind` attribute. If the response payload contains a `kind` field that matches the pattern
   `service#resource` where `service` is ignored and `resource` is equal to the name of a resource
   type then the response is assumed to be a resource collection of that type.

2. `Content-Type` header. If the response `Content-Type` header matches the pattern
   `vnd.*.resource+suffix;parameters` where `resource` is equal to the name of a resource type then
   the response is assumed to be a resource collection of that type.

3. `href` pattern. If the response payload contains a `href` field then it is matched against all
   the namespace resource type href patterns. If there are matches the response is assumed to be a
   resource collection of the type with the longest pattern that matches.

The algorithms tries these steps in order and if none return a match then an error is returned.



