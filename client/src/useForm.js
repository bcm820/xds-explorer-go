import { useState } from "react";

const useForm = (cb, initialState) => {
  const [values, setValues] = useState(initialState);

  const handleSubmit = event => {
    event.preventDefault();
    cb(values);
  };

  const handleChange = event => {
    event.persist();
    setValues(values => ({
      ...values,
      [event.target.name]: event.target.value
    }));
  };

  return { handleChange, handleSubmit, values };
};

export default useForm;
