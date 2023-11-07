"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
Object.defineProperty(exports, "__esModule", { value: true });
const galleta = __importStar(require("../dist/index"));
const sampleCookie = "eyJTdGF0ZVN0cmluZyI6IlRlc3RWYWx1ZSIsIlNpZyI6eyJGb3JtYXQiOiJzc2gtcnNhIiwiQmxvYiI6IkNCUHdES015RThUa3hRdm90d2xCc0ZxUUx3SG53cU1IMnNCU29jTXBzRVBUbFJTMXpOZ3hvNlloTjBTQVgxL1FaRGxoblFqbDJTOXluc21MRWJFa0NHTkJpM2lMN1FKK0F4NFZtSzZuTnBXcU1wSTd6cGFIZWVMY05ZcCt5RGFRd1dFT2pjbWk1d0ExdURvN3lLclZoWlFQamVTRHVBK1ZHM0JaSjZ3OWIwb0M0di9sVDFpMDgzSmZHUWpCWW5pS1lMaGZDeC9zQko4T2xuQWgxSDZmdi9nZU1MbWhuWVpJcktDdC84K01BSm5YTWM2UEZ3TlFqVHhWOVoySUZBcDc5S05VYzl0YmtxcjRwRzFTaWszbW9rczJES1U5SjdnYVlITGxsRFB0bloxaWJKeDdUY290S2FJQmRrTHBnTkJMSkNJRjZpWnc1TGNoTlAvS1Z2MktQQT09IiwiUmVzdCI6bnVsbH19";
describe("GetCookie Tests", () => {
    it("should return a json object", () => {
        console.log(galleta.decodeCookie(sampleCookie));
        true;
    });
});
//# sourceMappingURL=index.spec.js.map