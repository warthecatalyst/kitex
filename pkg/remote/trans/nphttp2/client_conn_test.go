/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nphttp2

import (
	"testing"

	"github.com/cloudwego/kitex/internal/test"
)

func TestClientConn(t *testing.T) {
	// init
	connPool := newMockConnPool()
	ctx := newMockCtxWithRPCInfo()
	conn, err := connPool.Get(ctx, "tcp", mockAddr0, newMockConnOption())
	test.Assert(t, err == nil, err)
	defer conn.Close()
	clientConn := conn.(*clientConn)

	// test LocalAddr()
	la := clientConn.LocalAddr()
	test.Assert(t, la == nil)

	// test RemoteAddr()
	ra := clientConn.RemoteAddr()
	test.Assert(t, ra.String() == mockAddr0)

	// test Read()
	testStr := "1234567890"
	mockStreamRecv(clientConn.s, testStr)
	testByte := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	n, err := clientConn.Read(testByte)
	test.Assert(t, err == nil, err)
	test.Assert(t, n == len(testStr))
	test.Assert(t, string(testByte) == testStr)

	// test write()
	testByte = []byte(testStr)
	n, err = clientConn.Write(testByte)
	test.Assert(t, err == nil, err)
	test.Assert(t, n == len(testByte))

	// test write() short write error
	n, err = clientConn.Write([]byte(""))
	test.Assert(t, err != nil, err)
	test.Assert(t, n == 0)
}

//func TestClientConnRead(t *testing.T) {
//	mockey.PatchConvey("TestClientConnRead", t, func() {
//		mockey.Mock((*grpc.Stream).Read).Return(1, io.EOF).Build()
//		mockey.Mock((*grpc.Stream).Status).Return(status.New(codes.Internal, "not found")).Build()
//		mockey.Mock((*grpc.Stream).BizStatusErr).Return(kerrors.NewGRPCBizStatusError(404, "not found")).Build()
//		s := &grpc.Stream{}
//		cli := &clientConn{s: s}
//		n, err := cli.Read(nil)
//		test.Assert(t, n == 1)
//		bizErr, _ := kerrors.FromBizStatusError(err)
//		test.Assert(t, bizErr.BizStatusCode() == 404)
//		test.Assert(t, bizErr.BizMessage() == "not found")
//	})
//	mockey.PatchConvey("TestClientConnRead", t, func() {
//		mockey.Mock((*grpc.Stream).Read).Return(1, io.EOF).Build()
//		mockey.Mock((*grpc.Stream).Status).Return(status.New(codes.Internal, "not found")).Build()
//		mockey.Mock((*grpc.Stream).BizStatusErr).Return(nil).Build()
//		s := &grpc.Stream{}
//		cli := &clientConn{s: s}
//		n, err := cli.Read(nil)
//		test.Assert(t, err != nil)
//		_, isBizErr := kerrors.FromBizStatusError(err)
//		test.Assert(t, !isBizErr)
//		test.Assert(t, n == 1)
//	})
//}
