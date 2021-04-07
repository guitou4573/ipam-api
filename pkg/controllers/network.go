package controllers

import (
	"log"
  "github.com/guitou4573/ipam/pkg/stores"
  "github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
)

type NetworkController struct{
	SubnetStore *stores.SubnetStore
  VPCStore *stores.VPCStore
}

type NetworkControllerOptions struct {
	SubnetStore *stores.SubnetStore
  VPCStore *stores.VPCStore
}

func NewNetworkController(options NetworkControllerOptions) *NetworkController {
	return &NetworkController{
		SubnetStore: options.SubnetStore,
    VPCStore: options.VPCStore,
	}
}

// /v1/vpc/:vpcid/subnet
// @Summary Get multiple subnets
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/vpc/:vpcid/subnet [get]
func (m NetworkController) ListSubnets(c *gin.Context) {
	tracer := opentracing.GlobalTracer()
	parentSpan, _ := c.Get("tracing-req-context")
	span := tracer.StartSpan(
		"subnet-get-by-vpc",
		opentracing.ChildOf(parentSpan.(opentracing.Span).Context()),
	)
  defer span.Finish()
  vpcId := c.Param("vpcid")
	subnets, err := m.SubnetStore.GetByVPCId(
    vpcId,
  )
	if err != nil {
    log.Println(err)
		c.JSON(500, gin.H{
			"error": "something happened",
		})
    return
	}

	c.JSON(200, gin.H{
		"subnets": subnets,
	})
}

// /v1/vpc
// @Summary Get multiple subnets
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /v1/vpc [get]
func (m NetworkController) ListVPCs(c *gin.Context) {
  tracer := opentracing.GlobalTracer()
	parentSpan, _ := c.Get("tracing-req-context")
	span := tracer.StartSpan(
		"get-vpc",
		opentracing.ChildOf(parentSpan.(opentracing.Span).Context()),
	)
  defer span.Finish()
  section := c.Param("section")
  offset := 0
  limit := 10
  vpcs, err := m.VPCStore.GetBySection(section, offset, limit)
  if err != nil {
    log.Println(err)
		c.JSON(500, gin.H{
			"error": "something happened",
		})
    return
	}

	c.JSON(200, gin.H{
		"vpcs": vpcs,
	})
}
