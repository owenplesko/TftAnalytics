import { Button } from "@/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { getUnitStats } from "@/services/getUnitStats";
import { UnitStat } from "@/services/types";
import { createFileRoute } from "@tanstack/react-router";
import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  getSortedRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { ArrowUpDown } from "lucide-react";

export const Route = createFileRoute("/units")({
  component: UnitStatsPage,
  loader: async ({ context: { queryClient } }) => {
    const unitStats = await queryClient.ensureQueryData(getUnitStats());

    return { unitStats };
  },
});

function UnitStatsPage() {
  return <UnitStatsTable />;
}

const columnHelper = createColumnHelper<UnitStat>();

const columns = [
  columnHelper.accessor("unit", {
    header: "Unit",
    cell: (props) => props.getValue(),
    enableSorting: false,
  }),
  columnHelper.accessor("avgPlacement", {
    header: "Avg Place",
    cell: (props) => props.getValue().toPrecision(3),
  }),
  columnHelper.accessor("frequency", {
    header: "Frequency",
    cell: (props) => props.getValue(),
  }),
];

function UnitStatsTable() {
  const { unitStats } = Route.useLoaderData();
  const table = useReactTable({
    columns,
    data: unitStats,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
  });

  return (
    <div className="w-[1000px] space-y-2">
      <h1 className="col-span-2 text-4xl font-semibold">TFT Unit Stats</h1>
      <div className="grid w-64 grid-cols-[auto_auto] rounded-md bg-card p-4 text-xs">
        <span>Last Updated: </span>
        <span className="text-right">a few seconds ago</span>
        <span>Comps Analyzed: </span>
        <span className="text-right">69,420</span>
      </div>
      <div className="w-full rounded-md border bg-card">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow>
                {headerGroup.headers.map((header) => (
                  <TableHead className="text-center first:text-left last:text-right">
                    {header.column.getCanSort() ? (
                      <Button
                        variant="ghost"
                        onClick={header.column.getToggleSortingHandler()}
                      >
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                        <ArrowUpDown size={16} />
                      </Button>
                    ) : (
                      <div>
                        {flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                      </div>
                    )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell
                    key={cell.id}
                    className="text-center first:text-left last:text-right"
                  >
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
