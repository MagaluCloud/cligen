# Confirmation

## Simple Ask

This type of confirmation is used when you need to ask the user for a simple confirmation.
### Example
```
confirmation:
  enabled: true
  type: simple-ask
  text: Are you sure you want to reset the container registry's password?
```

## Just Type

This type of confirmation is used when you need to ask the user for a specific value to confirm the action.
### Example
```
confirmation:
  enabled: true
  type: just-type
  value: alohomora
  text: "Digite `alohomora` para continuar:"
```

## Type Value

This type of confirmation is used when you need to ask the user for a specific field value to confirm the action.
### Example
```

confirmation:
  enabled: true
  type: type-value
  field: id
  text: "Informe '%s' para confirmar:"                
```