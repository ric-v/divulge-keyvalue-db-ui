import * as React from "react";
import { DataGrid } from "@mui/x-data-grid";
import { randomInt, randomUserName } from "@mui/x-data-grid-generator";
import Box from "@mui/material/Box";
import Stack from "@mui/material/Stack";
import Button from "@mui/material/Button";

const columns = [
  { field: "id" },
  { field: "username", width: 150 },
  { field: "age", width: 80, type: "number" },
];

let idCounter = 0;
const createRandomRow = () => {
  idCounter += 1;
  return { id: idCounter, username: randomUserName(), age: randomInt(10, 80) };
};

export default function UpdateRowsProp() {
  const [rows, setRows] = React.useState(() => [
    createRandomRow(),
    createRandomRow(),
    createRandomRow(),
    createRandomRow(),
  ]);

  const handleUpdateRow = () => {
    setRows((prevRows) => {
      const rowToUpdateIndex = randomInt(0, rows.length - 1);

      return prevRows.map((row, index) =>
        index === rowToUpdateIndex
          ? { ...row, username: randomUserName() }
          : row
      );
    });
  };

  const handleUpdateAllRows = () => {
    setRows(rows.map((row) => ({ ...row, username: randomUserName() })));
  };

  const handleDeleteRow = () => {
    setRows((prevRows) => {
      const rowToDeleteIndex = randomInt(0, prevRows.length - 1);
      return [
        ...rows.slice(0, rowToDeleteIndex),
        ...rows.slice(rowToDeleteIndex + 1),
      ];
    });
  };

  const handleAddRow = () => {
    setRows((prevRows) => [...prevRows, createRandomRow()]);
  };

  return (
    <div style={{ width: "100%" }}>
      <Stack
        sx={{ width: "100%", mb: 1 }}
        direction="row"
        alignItems="flex-start"
        columnGap={1}
      >
        <Button size="small" onClick={handleUpdateRow}>
          Update a row
        </Button>
        <Button size="small" onClick={handleUpdateAllRows}>
          Update all rows
        </Button>
        <Button size="small" onClick={handleDeleteRow}>
          Delete a row
        </Button>
        <Button size="small" onClick={handleAddRow}>
          Add a row
        </Button>
      </Stack>
      <Box sx={{ height: 400, bgcolor: "background.paper" }}>
        <DataGrid hideFooter rows={rows} columns={columns} />
      </Box>
    </div>
  );
}
