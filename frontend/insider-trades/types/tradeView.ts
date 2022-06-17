export type TradeView = {
  ID: string;
  SecFilingsID: string;
  TransactionTypeName: string;
  AveragePrice: number;
  TotalShares: number;
  TotalValue: number;
  CreatedAt: protoTypestamp;
  URL: string;
  InsiderID: string;
  CompanyID: string;
  OfficerPosition: string;
  ReportedOn: string;
  InsiderCik: number;
  InsiderName: string;
  CompanyCik: number;
  CompanyName: string;
  CompanyTicker: string;
};

export type protoTypestamp = {
  seconds: number;
  nanos: number;
};
