package main

import("github.com/gin-gonic/gin"
       "shef-boutique/db"
       "shef-boutique/product"
       "shef-boutique/category"
       "shef-boutique/subcategory"
       "shef-boutique/admin"
       "shef-boutique/stats"
       "shef-boutique/search"
       "shef-boutique/promotions"
      )

func main(){

	r := gin.Default()
	db.Connect()

	r.GET("/product/category/:id",product.AllByCategory)
	r.GET("/product/:id",product.GetOne)
	r.POST("/product",product.Add)
	r.PUT("/product/:id",product.Update)
	r.PUT("/product/price",product.ToggleHidePrice)
	r.GET("/product/price",product.HidePrice)
	r.PUT("/product/:id/stock",product.ToggleStock)
	r.PUT("/product/:id/newarrival",product.ToggleNewArrival)
	r.DELETE("/product/:id",product.Delete)

	r.GET("/category",category.All)
	r.POST("/category",category.Add)
	r.DELETE("/category/:id",category.Delete)

	r.GET("/subcategory/category/:id",subcategory.All)
	r.POST("/subcategory",subcategory.Add)
	r.DELETE("/subcategory/:id",subcategory.Delete)

	r.GET("/promotion/:id",promotions.GetOne)
	r.POST("/promotion",promotions.Add)
	r.GET("/promotion",promotions.All)
	r.DELETE("/promotion/:id",promotions.Delete)

	r.POST("/admin/register",admin.Register)
	r.POST("/admin/login",admin.Login)
	r.POST("/admin/contact",admin.UpdateContact)
	r.GET("/admin/contact",admin.GetContact)

	r.POST("/stats/:id",stats.StatsByDay)
	r.GET("/search",search.Products)
	r.Run();
}
