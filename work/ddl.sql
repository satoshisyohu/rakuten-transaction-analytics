
CREATE TABLE `saichan-435607.transaction_analtics.Transactions` (
                                                                    TransactionDate DATE   NOT NULL,  -- REQUIRED
                                                                    Detail          STRING NOT NULL,  -- REQUIRED
                                                                    Amount          INT64,            -- NULLABLE
                                                                    Category        STRING,            -- NULLABLE
                                                                    TransactionType STRING,            -- REQUIRED
                                                                    Id STRING NOT NULL,  -- REQUIRED
);



CREATE TABLE `saichan-435607.transaction_analtics.TransactionReports` (
                                                                  Id             STRING   NOT NULL,   -- REQUIRED
                                                                  YearMonth      STRING   NOT NULL,   -- REQUIRED
                                                                  BaseAmounts    INT64    NOT NULL,   -- REQUIRED
                                                                  TotalAmount    INT64    NOT NULL,   -- REQUIRED
                                                                  FoodExpenses   INT64    NOT NULL,   -- REQUIRED
                                                                  WasteExpenses  INT64    NOT NULL,   -- REQUIRED
                                                                  OtherExpenses  INT64    NOT NULL,   -- REQUIRED
                                                                  FixedCosts     INT64    NOT NULL,   -- REQUIRED
                                                                  VariableCosts  INT64    NOT NULL,   -- REQUIRED
                                                                  Savings        INT64    NOT NULL,  -- REQUIRED
                                                                  Score          FLOAT64  NOT NULL   -- REQUIRED
);