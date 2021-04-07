PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS `token` (
  "token" VARCHAR(255) PRIMARY KEY NOT NULL,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS `subnet` (
  "netaddr" INTEGER CHECK("netaddr">=0),
  "mask" INTEGER CHECK("mask">=0),
  "idvpc" VARCHAR(255) NOT NULL,
  "description" VARCHAR(255),
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY ("netaddr","mask")
  CONSTRAINT "fk_subnet_vpc"
    FOREIGN KEY ("idvpc")
    REFERENCES "vpc"("id")
);
CREATE INDEX IF NOT EXISTS "subnet.fk_subnet_vpc_idx" ON "subnet" ("idvpc");

CREATE TABLE IF NOT EXISTS `vpc`(
  "id" VARCHAR(255) NOT NULL,
  "idsection" VARCHAR(255) NOT NULL,
  "description" VARCHAR(255),
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY("id")
  CONSTRAINT "fk_vpc_section"
    FOREIGN KEY("idsection")
    REFERENCES "section"("id")
);
CREATE INDEX IF NOT EXISTS "vpc.fk_vpc_section_idx" ON "vpc" ("idsection");

CREATE TABLE IF NOT EXISTS `section` (
  "id" VARCHAR(255) NOT NULL,
  "parent" VARCHAR(255),
  "description" VARCHAR(255),
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY("id")
  CONSTRAINT "fk_section_section"
    FOREIGN KEY("parent")
    REFERENCES "section"("id")
);
CREATE INDEX IF NOT EXISTS "section.fk_section_section_idx" ON "section" ("parent");

INSERT INTO `section` ("id", "description") VALUES ("aws-dev", "AWS dev account");
INSERT INTO `vpc` ("id", "idsection", "description") VALUES ("vpc-349234", "aws-dev", "demo-vpc");
INSERT INTO `subnet` ("netaddr", "mask", "idvpc", "description") VALUES (inet_aton("10.0.0.0"), inet_aton("255.255.255.0"), "vpc-349234", "Test subnet");
