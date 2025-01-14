package handlers

import (
	"encoding/json"
	"net"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/services"
	"ykc-proxy-server/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func StartChargingRouter(c *gin.Context) {
	clientID := c.DefaultQuery("clientID", "")
	if clientID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}
	err := services.StartCharging(clientID)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"status": "message sent"})
}

func StopChargingRouter(c *gin.Context) {

	clientID := c.DefaultQuery("clientID", "")
	if clientID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}
	err := services.StopCharging(clientID)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"status": "message sent"})

}

func VerificationResponseRouter(c *gin.Context) {
	var req dtos.VerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := services.ResponseToVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func VerificationHandler(opt *dtos.Options, buf []byte, hex []string, header *dtos.Header, conn net.Conn,
) {
	msg := services.Verification(opt, buf, hex, header, conn)

	if opt.AutoVerification {
		m := &dtos.VerificationResponseMessage{
			Header: &dtos.Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:     msg.Id,
			Result: true,
		}
		_ = services.ResponseToVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("01", b)
	}
}

func HeartbeatHandler(buf []byte, header *dtos.Header, conn net.Conn) {
	_ = services.Hearthbeat(buf, header, conn)
	// Send Heartbeat Response
	_ = services.SendHeartbeatResponse(conn, header)
}

func BillingModelVerificationHandler(opt *dtos.Options, hex []string, header *dtos.Header, conn net.Conn) {
	msg := services.BillingModelVerification(opt, hex, header, conn)
	//auto response
	if opt.AutoBillingModelVerify {
		m := &dtos.BillingModelVerificationResponseMessage{
			Header: &dtos.Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:               msg.Id,
			BillingModelCode: msg.BillingModelCode,
			Result:           true,
		}
		_ = services.ResponseToBillingModelVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("05", b)
	}
}

func BillingModelRequestMessageHandler(opt *dtos.Options, hex []string, header *dtos.Header, conn net.Conn) {
	msg := services.BillingModelRequestMessage(opt, hex, header, conn)
	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("09", b)
	}
}

func BillingModelResponseMessageHandler(c *gin.Context) {
	var req dtos.BillingModelResponseMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendBillingModelResponseMessage(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func BillingModelVerificationResponseHandler(c *gin.Context) {
	var req dtos.BillingModelVerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := services.ResponseToBillingModelVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteBootstrapRequestHandler(c *gin.Context) {
	var req dtos.RemoteBootstrapRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendRemoteBootstrapRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteShutdownRequestHandler(c *gin.Context) {
	var req dtos.RemoteShutdownRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendRemoteShutdownRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func TransactionRecordConfirmedHandler(c *gin.Context) {
	var req dtos.TransactionRecordConfirmedMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendTransactionRecordConfirmed(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SetBillingModelRequestHandler(c *gin.Context) {
	var req dtos.SetBillingModelRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendSetBillingModelRequestMessage(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteRebootRequestMessageHandler(c *gin.Context) {
	var req dtos.RemoteRebootRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendRemoteRebootRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SubmitFinalStatusHandler(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn,
) {
	data := services.SubmitFinalStatus(opt, buf, header, conn)
	err := utils.SendMessage(conn, data)
	if err != nil {
		log.Errorf("Failed to send Submit Final Status response: %v", err)
	} else {
		log.Debug("Sent Submit Final Status response successfully")
	}
}

func DeviceLoginHandler(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn,
) {
	msg, resMessage := services.DeviceLogin(opt, buf, header, conn)
	// Forward the Device Login message to an external system (optional)
	if msg != nil {
		if opt.MessageForwarder != nil {
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Errorf("Failed to marshal Device Login message: %v", err)
				return
			}
			err = opt.MessageForwarder.Publish("81", jsonMsg)
			if err != nil {
				log.Errorf("Failed to publish Device Login message: %v", err)
			}
		}
	}
	if resMessage != nil {
		utils.SendMessage(conn, resMessage)
	}
}

func RemoteStartHandler(buf []byte, header *dtos.Header, conn net.Conn) {
	services.RemoteStart(buf, header, conn)
	// if data != nil {
	// 	err := utils.SendMessage(conn, data)
	// 	if err != nil {
	// 		log.Errorf("Failed to send Remote Start response: %v", err)
	// 	} else {
	// 		log.Debug("Sent Remote Start response successfully")
	// 	}
	// }

}

func RemoteStopHandler(buf []byte, header *dtos.Header, conn net.Conn) {
	data := services.RemoteStop(buf, header, conn)
	if data != nil {
		err := utils.SendMessage(conn, data)
		if err != nil {
			log.Errorf("Failed to send Remote Start response: %v", err)
		} else {
			log.Debug("Sent Remote Start response successfully")
		}
	}
}

func ChargingPortDataHandler(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn) {
	// Parse the CMD088 message
	msg := services.ChargingPortData(opt, buf, header, conn)

	// Forward the Charging Port Data message to an external system (optional)
	if msg != nil {
		if opt.MessageForwarder != nil {
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Errorf("Failed to marshal Charging Port Data message: %v", err)
				return
			}
			err = opt.MessageForwarder.Publish("88", jsonMsg)
			if err != nil {
				log.Errorf("Failed to publish Charging Port Data message: %v", err)
			}
		}
	}

}
