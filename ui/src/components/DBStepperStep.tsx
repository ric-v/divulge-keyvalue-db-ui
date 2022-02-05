import { StepLabel, StepContent, Box, Button, Typography } from "@mui/material";

type Props = {
  stepTitle: string;
  buttonName: string;
  nextBtnDisabled: boolean;
  backBtnDisabled: boolean;
  handleNext: () => void;
  handleBack: () => void;
  stepContent: () => JSX.Element;
};

const DBStepperStep = ({
  stepTitle,
  buttonName,
  nextBtnDisabled,
  backBtnDisabled,
  handleNext,
  handleBack,
  stepContent,
}: Props) => {
  return (
    <>
      <StepLabel>
        <Typography variant="h3" color="primary.light">
          {stepTitle}
        </Typography>
      </StepLabel>
      <StepContent>
        {stepContent()}
        <Box sx={{ mb: 2 }}>
          <div>
            <Button
              disabled={nextBtnDisabled}
              variant="contained"
              onClick={handleNext}
              sx={{ mt: 1, mr: 1 }}
            >
              {buttonName}
            </Button>
            <Button
              disabled={backBtnDisabled}
              onClick={handleBack}
              sx={{ mt: 1, mr: 1 }}
            >
              Back
            </Button>
          </div>
        </Box>
      </StepContent>
    </>
  );
};

export default DBStepperStep;
