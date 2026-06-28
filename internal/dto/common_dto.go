package dto 

type StandardResponse struct { 
	Success bool        `json:"success"`          
	Message string      `json:"message"`          
	Data    interface{} `json:"data,omitempty"`   
	Errors  interface{} `json:"errors,omitempty"` 
}
