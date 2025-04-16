# Common Instructions for Rust Porting

## General Guidelines

- Follow Rust best practices and idioms
- Use proper error handling with Result types
- Implement proper documentation
- Add TODOs for incomplete functionality with detailed explanations.
- document the translation by writing a line in a file in the form of "source" -> "destination"
- add third-party libraries as needed.
- use htmx for the frontend.

## Code Structure

- Organize code into modules
- Use appropriate visibility modifiers
- Follow the Rust module system conventions
- All Rust code should go into the @rust_app directory
- integrate the file(s) into the existing appliation

## Testing

- Write unit tests for all public functions
- Include integration tests where appropriate
- Use property-based testing for complex logic

## Performance Considerations

- Minimize allocations
- Use appropriate data structures
- Consider async/await for I/O operations

## Security

- Validate all input
- Use proper authentication and authorization
- Follow the principle of least privilege

## Prompt Output

- dont provide additional explanatory information.
