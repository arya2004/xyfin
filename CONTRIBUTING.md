# Contributing to Xyfin

We are excited that you're interested in contributing to Xyfin! Whether you're fixing a bug, adding new features, or improving the documentation, your contributions are highly appreciated. This guide will help you get started.

## Getting Started

### 1. Fork the Repository

Start by forking the Xyfin repository. This will create a copy of the repository under your GitHub account.

### 2. Clone the Repository

Clone the forked repository to your local machine:

```bash
git clone https://github.com/arya2004/xyfin.git
cd xyfin
```

### 3. Set Up the Environment

Follow the instructions in the `README.md` file to set up your development environment. Ensure you have Docker, Go, and all necessary dependencies installed.

### 4. Create a Branch

Create a new branch for your work. Use a descriptive name that reflects the scope of your contribution:

```bash
git checkout -b feature/your-feature-name
```

### 5. Make Your Changes

Make your changes, whether it’s fixing a bug, adding a feature, or improving documentation. Be sure to:

- Write clean, readable code.
- Follow the project's coding style and conventions.
- Include comments where necessary.
- Write or update tests for any new functionality.

### 6. Run Tests

Ensure all tests pass before committing your changes:

```bash
make test
```

### 7. Commit Your Changes

Once you're satisfied with your changes, commit them to your branch:

```bash
git add .
git commit -m "Description of your changes"
```

### 8. Push Your Changes

Push your changes to your forked repository:

```bash
git push origin feature/your-feature-name
```

### 9. Create a Pull Request

Go to the original Xyfin repository on GitHub and create a new Pull Request from your branch. Provide a clear description of your changes and the problem they solve.

## Code Review Process

Your Pull Request will be reviewed by other contributors or maintainers. They may ask for changes or clarifications before it can be merged. Please be responsive to feedback and make any necessary updates.

## Contribution Guidelines

### Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Keep the codebase clean and maintainable.
- Write meaningful commit messages.

### Testing

- Ensure that your code changes do not break existing tests.
- Write tests for any new functionality you introduce.
- Aim for high test coverage.

### Documentation

- Update documentation to reflect your changes, especially if they affect the public API.
- Ensure your documentation is clear and easy to understand.

### Merging

- Once your Pull Request is approved, a maintainer will merge it into the main branch.
- Make sure your branch is up-to-date with the main branch before merging.

## Reporting Issues

If you find a bug or have a feature request, please create an issue on GitHub. Provide as much detail as possible, including steps to reproduce the issue or a detailed description of the feature you would like to see.

## Code of Conduct

Please note that this project is governed by a [Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project, you agree to abide by its terms.

## Thank You!

Thank you for contributing to Xyfin! Your help is greatly appreciated, and we're thrilled to have you on board. If you have any questions, feel free to ask in an issue or reach out to the maintainers.

Happy coding!