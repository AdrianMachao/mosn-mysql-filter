package mysql

/*
*   https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_packets.html#sect_protocol_basic_packets_sending_mt_16mb.

* 0           1           2           3           4           5           6           7           8
* +-----------------------------------+-----------+-----------------------------------------------+
* |             payload_length        |sequence_id|                      payload                  |
* +-----------------------------------+-----------+-----------------------------------------------+

 */

// http://dev.mysql.com/doc/internals/en/client-server-protocol.html

// Read packet to buffer 'data'
//func (mc *mysqlConn) readPacket() ([]byte, error) {
//	var prevData []byte
//	for {
//		// read packet header
//		data, err := mc.buf.readNext(4)
//		if err != nil {
//			if cerr := mc.canceled.Value(); cerr != nil {
//				return nil, cerr
//			}
//			errLog.Print(err)
//			mc.Close()
//			return nil, ErrInvalidConn
//		}
//
//		// packet length [24 bit]
//		pktLen := int(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16)
//
//		// check packet sync [8 bit]
//		if data[3] != mc.sequence {
//			if data[3] > mc.sequence {
//				return nil, ErrPktSyncMul
//			}
//			return nil, ErrPktSync
//		}
//		mc.sequence++
//
//		// packets with length 0 terminate a previous packet which is a
//		// multiple of (2^24)-1 bytes long
//		if pktLen == 0 {
//			// there was no previous packet
//			if prevData == nil {
//				errLog.Print(ErrMalformPkt)
//				mc.Close()
//				return nil, ErrInvalidConn
//			}
//
//			return prevData, nil
//		}
//
//		// read packet body [pktLen bytes]
//		data, err = mc.buf.readNext(pktLen)
//		if err != nil {
//			if cerr := mc.canceled.Value(); cerr != nil {
//				return nil, cerr
//			}
//			errLog.Print(err)
//			mc.Close()
//			return nil, ErrInvalidConn
//		}
//
//		// return data if this was the last packet
//		if pktLen < maxPacketSize {
//			// zero allocations for non-split packets
//			if prevData == nil {
//				return data, nil
//			}
//
//			return append(prevData, data...), nil
//		}
//
//		prevData = append(prevData, data...)
//	}
//}
//
//// Write packet buffer 'data'
//func (mc *mysqlConn) writePacket(data []byte) error {
//	pktLen := len(data) - 4
//
//	if pktLen > mc.maxAllowedPacket {
//		return ErrPktTooLarge
//	}
//
//	// Perform a stale connection check. We only perform this check for
//	// the first query on a connection that has been checked out of the
//	// connection pool: a fresh connection from the pool is more likely
//	// to be stale, and it has not performed any previous writes that
//	// could cause data corruption, so it's safe to return ErrBadConn
//	// if the check fails.
//	if mc.reset {
//		mc.reset = false
//		conn := mc.netConn
//		if mc.rawConn != nil {
//			conn = mc.rawConn
//		}
//		var err error
//		if mc.cfg.CheckConnLiveness {
//			if mc.cfg.ReadTimeout != 0 {
//				err = conn.SetReadDeadline(time.Now().Add(mc.cfg.ReadTimeout))
//			}
//			if err == nil {
//				err = connCheck(conn)
//			}
//		}
//		if err != nil {
//			errLog.Print("closing bad idle connection: ", err)
//			mc.Close()
//			return driver.ErrBadConn
//		}
//	}
//
//	for {
//		var size int
//		if pktLen >= maxPacketSize {
//			data[0] = 0xff
//			data[1] = 0xff
//			data[2] = 0xff
//			size = maxPacketSize
//		} else {
//			data[0] = byte(pktLen)
//			data[1] = byte(pktLen >> 8)
//			data[2] = byte(pktLen >> 16)
//			size = pktLen
//		}
//		data[3] = mc.sequence
//
//		// Write packet
//		if mc.writeTimeout > 0 {
//			if err := mc.netConn.SetWriteDeadline(time.Now().Add(mc.writeTimeout)); err != nil {
//				return err
//			}
//		}
//
//		n, err := mc.netConn.Write(data[:4+size])
//		if err == nil && n == 4+size {
//			mc.sequence++
//			if size != maxPacketSize {
//				return nil
//			}
//			pktLen -= size
//			data = data[size:]
//			continue
//		}
//
//		// Handle error
//		if err == nil { // n != len(data)
//			mc.cleanup()
//			errLog.Print(ErrMalformPkt)
//		} else {
//			if cerr := mc.canceled.Value(); cerr != nil {
//				return cerr
//			}
//			if n == 0 && pktLen == len(data)-4 {
//				// only for the first loop iteration when nothing was written yet
//				return errBadConnNoWrite
//			}
//			mc.cleanup()
//			errLog.Print(err)
//		}
//		return ErrInvalidConn
//	}
//}
