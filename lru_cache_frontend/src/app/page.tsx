"use client";
import { useEffect, useState, useRef } from "react";
import { DataGrid } from '@mui/x-data-grid';
import * as React from 'react';
import { styled } from '@mui/material/styles';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import axios from 'axios'

let data = {
  key: 1,
  value: 24,
  timestamp: Date.now()
}
interface Res {
  key: Number,
  value: Number,
  timestamp: Date

}
let baseUrl = 'http://localhost:5000/api/'
const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: theme.palette.common.black,
    color: theme.palette.common.white,
  },
  [`&.${tableCellClasses.body}`]: {
    fontSize: 14,
  },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
  '&:nth-of-type(odd)': {
    backgroundColor: theme.palette.action.hover,
  },
  // hide last border
  '&:last-child td, &:last-child th': {
    border: 0,
  },
}));

export default function Home() {
  const [data, setData] = useState<Res[] | []>();
  const keyRef = useRef<HTMLInputElement>(null);
  const valueRef = useRef<HTMLInputElement>(null);
  const getKeyRef = useRef<HTMLInputElement>(null);
  const [getData, setGetData] = useState<Res>()

  useEffect(() => {
    let interval = setInterval(async () => {
      fetch(baseUrl + `all/`).then(
        response => response.json()
      ).then(json => setData(json))

    }, 1000);
    return () => {
      clearInterval(interval);
    };
  }, []);

  return (
    <div className="h-fit flex justify-center">

      <div className="w-8/12 flex flex-col  items-center">
        <div className="bg-orange-400 flex h-16 w-4/5 font-bold rounded-md">
          <h1 className="mt-3 ml-3 font-black text-3xl">LRU-CACHE-UI</h1>
        </div>
        <div className="flex w-4/5 justify-evenly mt-16 mb-9">
          <div className=" bg-slate-200 w-2/5 flex flex-col p-5 rounded-md">
            <h3 className="font-semibold">GET Cache</h3>
            <label htmlFor="" className="font-mono m-2">Key</label>
            <input type="text" ref={getKeyRef} className="rounded-lg h-10" />
            <br></br>
            <button className="bg-orange-500 rounded-md p-2 text-base font-semibold" onClick={() => {
              axios.get(baseUrl + `${getKeyRef.current!.value}`).then((d) => {
                setGetData(d.data)
              })
            }}>
              Search
            </button>
            {getData?.key ?
              <div className="font-medium pt-6">
                <span>{`Key :  ${getData?.key}`}</span>
                <br />
                <span>{`Value :  ${getData?.value}`}</span>
                <br />
                <span>{`Time :  ${new Date(getData?.timestamp)}`}</span>
              </div>
              : null}
          </div>
          <div className=" bg-slate-200 w-2/5 flex flex-col p-5 rounded-md">
            <h3 className="font-semibold">SET Cache</h3>
            <div className="flex flex-col mt-2">
              <div className="flex flex-col mb-2">
                <label htmlFor="" className="font-mono ">Key</label>
                <input type="text" ref={keyRef} className=" rounded-lg  h-10" />
              </div>
              <div className="flex flex-col">
                <label htmlFor="" className="font-mono ">Value</label>
                <input type="text" ref={valueRef} className="rounded-lg  h-10" />

              </div>
            </div>
            <br />
            <button className="bg-orange-500 rounded-md p-2 text-base font-semibold" onClick={() => {
              console.log(keyRef.current!.value, valueRef.current!.value);
              axios.post(baseUrl, {
                key: parseInt(keyRef.current!.value),
                value: parseInt(valueRef.current!.value)
              })
            }}>
              Add
            </button>
          </div>
        </div>
        <div className="flex w-4/5 mt-10 r">
          <TableContainer component={Paper}>
            <Table sx={{ minWidth: 700 }} aria-label="customized table">
              <TableHead>
                <TableRow>
                  <StyledTableCell>KEY</StyledTableCell>
                  <StyledTableCell align="right">VALUE</StyledTableCell>
                  <StyledTableCell align="right">TIMESTAMP</StyledTableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {data?.map((row, index) => (
                  <StyledTableRow key={index}>
                    <StyledTableCell component="th" scope="row">
                      {`${row.key}`}
                    </StyledTableCell>
                    <StyledTableCell align="right">{`${row.value}`}</StyledTableCell>
                    <StyledTableCell align="right">{`${row.timestamp}`}</StyledTableCell>
                  </StyledTableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </div>
      </div>
    </div>
  );
}
