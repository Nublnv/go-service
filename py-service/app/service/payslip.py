from payslip.v1.payslip_pb2_grpc import DocumentServiceServicer
from payslip.v1.payslip_pb2 import GetPayslipRequest, GeneratePayslipResponse, GetPayslipResponse, GetPayslipResponse
from grpc.aio import ServicerContext
from motor.motor_asyncio import AsyncIOMotorClient

class DocumentServer(DocumentServiceServicer):

    mongo: AsyncIOMotorClient

    def __init__(self, mongo: AsyncIOMotorClient):
        self.mongo = mongo
        super().__init__()

    async def GeneratePayslip(self, request: GetPayslipRequest, context: ServicerContext) -> GeneratePayslipResponse:
        """Method for generate payslip in Python service and save it in mongo. Returns payslipid
        """
        ...
    

    async def GetPayslip(self, request: GetPayslipResponse, context: ServicerContext) -> GeneratePayslipResponse:
        """Method returns payslip document by id
        """
        ...