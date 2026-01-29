-- Migration: Create products table for SQL Server
IF NOT EXISTS (SELECT * FROM sys.objects WHERE object_id = OBJECT_ID(N'[dbo].[products]') AND type in (N'U'))
BEGIN
    CREATE TABLE [dbo].[products] (
        [id] NVARCHAR(36) NOT NULL PRIMARY KEY,
        [name] NVARCHAR(100) NOT NULL UNIQUE,
        [description] NVARCHAR(MAX) NULL,
        [price] DECIMAL(10,2) NOT NULL CONSTRAINT chk_price CHECK (price > 0),
        [stock] INT NOT NULL CONSTRAINT chk_stock CHECK (stock >= 0),
        [category] NVARCHAR(50) NOT NULL,
        [active] BIT NOT NULL DEFAULT 1,
        [created_at] DATETIME2 NOT NULL DEFAULT (SYSUTCDATETIME()),
        [updated_at] DATETIME2 NOT NULL DEFAULT (SYSUTCDATETIME())
    );

    CREATE INDEX idx_products_category ON [dbo].[products]([category]);
    CREATE INDEX idx_products_active ON [dbo].[products]([active]);
    CREATE INDEX idx_products_created_at ON [dbo].[products]([created_at]);
END

-- Optional: trigger to update updated_at on row modification
IF OBJECT_ID('dbo.trg_products_updated_at', 'TR') IS NULL
BEGIN
    EXEC('CREATE TRIGGER dbo.trg_products_updated_at
    ON dbo.products
    AFTER UPDATE
    AS
    BEGIN
        SET NOCOUNT ON;
        UPDATE dbo.products
        SET updated_at = SYSUTCDATETIME()
        FROM inserted i
        WHERE dbo.products.id = i.id;
    END')
END
