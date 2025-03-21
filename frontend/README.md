# Expense Manager frontend

## Frontend Assignment
Your task is to create a **frontend application using either React, Vue, or Angular.
We are limiting the choices to the frameworks that we use at SPIRIA. 
The main task of the frontend application is to interact with the given API (backend). 

The backend is serving the necessary endpoints to manage expenses. 
Should you need to know more about the backend, please refer to the [backend README](../backend/README.md).
If you need an endpoint that is not available, feel free to create it in the backend.
The backend endpoints are protected and require a valid JWT authorization header.

## Requirements:
- **Authentication**: Provide a login form to authenticate users.
- **Registration**: Allow users to register a new account.
- **List Expenses**: Display a list of all expenses.
- **Create Expense**: Provide a form to create a new expense.
- **Update Expense**: Allow users to update an existing expense.
- **Delete Expense**: Enable users to delete an expense.
- **Create Category**: Allow users to create a new category.
- **Logout**: Provide a way for users to logout (_there is no backend endpoint for this_).
- **JWT Token**: Store the JWT token securely.
- **Error Handling**: Properly handle errors and provide user feedback.

## Authorization:
All endpoints require a valid JWT token in the Authorization header.

## Recommendations:
You are free to choose any modern JavaScript framework, but we highly recommend using React, Vue, or Angular.

## Notes:
Ensure proper error handling and user feedback.
Follow best practices for code organization and component structure.
Make any styling decisions you see fit.
Make sure to handle JWT token storage securely.
Where you find the requirement ambiguous, make an assumption and document it in the README file.

# Backend:
The backend for this application is already coded. 
Checkout the README [here](../backend/README.md) for more information. 
Remember to run the backend before running the frontend. Instructions are provided in the backend README.
There is no need to worry about deployment, we assume that it will be running locally on `http://localhost:8080`.
There is a swagger documentation available at `http://localhost:8080/api/swagger/index.html`.

**Note**: Run the application in docker for convenience, that way you will not need to install docker for the assignment, 
nor worry about the database setup and go modules. checkout the backend README for more information.


# Submission:
Fork this repository and create a new branch with your name.
Once you have completed the assignment, create a pull request to the main branch.

> Please update this README with instructions on how to run your application and 
any assumptions you made during development, including any additional features you implemented.

## Questions:
If you have any questions, feel free to reach out to us through the recruiter.

# IMPORTANT NOTE:
We are looking for a simple, clean, and functional solution. Therefore, don't spend too much time on styling or extra features. 
If you are unable to complete the assignment within the time frame, please submit what you have completed. 
We will evaluate your submission based on the quality of the code and the features implemented.