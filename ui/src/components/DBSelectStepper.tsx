import * as React from "react";
import Box from "@mui/material/Box";
import Stepper from "@mui/material/Stepper";
import Step from "@mui/material/Step";
import Button from "@mui/material/Button";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import {
  FormControl,
  FormControlLabel,
  FormGroup,
  Radio,
  RadioGroup,
  Switch,
} from "@mui/material";
import DBStepperStep from "./DBStepperStep";
import FileUpload from "../controllers/FileUpload";
import NewFile from "../controllers/NewFile";

type Props = {
  dbName: string;
  setDbname: React.Dispatch<React.SetStateAction<string>>;
  dbKey: string;
  setDbkey: React.Dispatch<React.SetStateAction<string>>;
  status: string;
  setStatus: React.Dispatch<React.SetStateAction<string>>;
  dbtype: string;
  setDbtype: React.Dispatch<React.SetStateAction<string>>;
  loadView: boolean;
  setLoadView: React.Dispatch<React.SetStateAction<boolean>>;
};

export default function VerticalLinearStepper(
  this: any,
  {
    dbName,
    setDbname,
    dbKey,
    setDbkey,
    setStatus,
    dbtype,
    setDbtype,
    setLoadView,
  }: Props
) {
  const [activeStep, setActiveStep] = React.useState(0);
  const [checked, setChecked] = React.useState(false);

  const handleNext = (finish: boolean) => {
    setActiveStep((prevActiveStep) => prevActiveStep + 1);
    if (finish) setLoadView(true);
  };

  const handleBack = () => {
    setActiveStep((prevActiveStep) => prevActiveStep - 1);
  };

  const handleReset = () => {
    setActiveStep(0);
  };

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setChecked(event.target.checked);
  };

  const selectDatabaseTypeStep = () => {
    return (
      <FormControl>
        <RadioGroup
          row
          aria-labelledby="demo-row-radio-buttons-group-label"
          name="row-radio-buttons-group"
        >
          <FormControlLabel
            value="buntdb"
            control={<Radio />}
            label="Bunt DB"
            onChange={(e) => {
              e.preventDefault();
              setDbtype("buntdb");
            }}
          />
          <FormControlLabel
            value="boltdb"
            control={<Radio />}
            label="Bolt DB"
            onChange={(e) => {
              e.preventDefault();
              setDbtype("boltdb");
            }}
            disabled
          />
        </RadioGroup>
      </FormControl>
    );
  };

  const selectDatabaseStep = () => {
    return (
      <>
        <FormGroup>
          <FormControlLabel
            control={
              <Switch
                checked={checked}
                onChange={handleChange}
                inputProps={{ "aria-label": "controlled" }}
                disabled
              />
            }
            label={checked ? "Create new database" : "Upload existing database"}
          />
        </FormGroup>
        {checked ? (
          <NewFile
            setDbkey={setDbkey}
            dbName={dbName}
            setDbname={setDbname}
            setStatus={setStatus}
          />
        ) : (
          <FileUpload
            setDbkey={setDbkey}
            setDbname={setDbname}
            setStatus={setStatus}
          />
        )}
      </>
    );
  };

  const loadViewStep = () => {
    return (
      <>
        <Typography variant="h5">{dbName} database loaded!</Typography>
        <Typography>
          Click finish setup to load the database grid view
        </Typography>
      </>
    );
  };

  return (
    <Box sx={{ maxWidth: "100%" }}>
      <Stepper activeStep={activeStep} orientation="vertical">
        <Step key="1">
          <DBStepperStep
            stepTitle={"Select Database Type"}
            buttonName={"Next"}
            nextBtnDisabled={dbtype ? false : true}
            backBtnDisabled={true}
            handleNext={handleNext.bind(this, false)}
            handleBack={handleBack}
            stepContent={selectDatabaseTypeStep}
          />
        </Step>

        <Step key="2">
          <DBStepperStep
            stepTitle={"Upload / Create New Database"}
            buttonName={"Next"}
            nextBtnDisabled={dbKey && dbName ? false : true}
            backBtnDisabled={false}
            handleNext={handleNext.bind(this, false)}
            handleBack={handleBack}
            stepContent={selectDatabaseStep}
          />
        </Step>

        <Step key="3">
          <DBStepperStep
            stepTitle={"Load Database"}
            buttonName={"Finish Setup"}
            nextBtnDisabled={dbKey && dbName ? false : true}
            backBtnDisabled={false}
            handleNext={handleNext.bind(this, true)}
            handleBack={handleBack}
            stepContent={loadViewStep}
          />
        </Step>
      </Stepper>
      {activeStep === 3 && (
        <Paper square elevation={0} sx={{ p: 3 }}>
          <Typography>All steps completed - you&apos;re finished</Typography>
          <Button onClick={handleReset} sx={{ mt: 1, mr: 1 }}>
            Reset
          </Button>
        </Paper>
      )}
    </Box>
  );
}
