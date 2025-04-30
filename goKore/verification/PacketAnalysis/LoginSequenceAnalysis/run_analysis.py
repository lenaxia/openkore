#!/usr/bin/env python3
"""
Run all login sequence analysis tools and generate a comprehensive report.
"""

import os
import subprocess
import sys
import time

def run_script(script_name):
    """Run a Python script and return its output."""
    print(f"Running {script_name}...")
    start_time = time.time()
    
    try:
        result = subprocess.run(
            [sys.executable, script_name],
            capture_output=True,
            text=True,
            check=True
        )
        print(f"✓ {script_name} completed successfully in {time.time() - start_time:.2f} seconds")
        return result.stdout
    except subprocess.CalledProcessError as e:
        print(f"✗ Error running {script_name}: {e}")
        print(f"Error output: {e.stderr}")
        return None

def generate_comprehensive_report(output_dir):
    """Generate a comprehensive report combining all analysis results."""
    report_parts = []
    
    # Add README content
    readme_path = os.path.join(output_dir, "README.md")
    if os.path.exists(readme_path):
        with open(readme_path, 'r') as f:
            report_parts.append(f.read())
    
    # Add login sequence summary
    summary_path = os.path.join(output_dir, "login_sequence_summary.md")
    if os.path.exists(summary_path):
        report_parts.append("\n\n---\n\n")
        with open(summary_path, 'r') as f:
            report_parts.append(f.read())
    
    # Add packet structures
    structures_path = os.path.join(output_dir, "packet_structures.md")
    if os.path.exists(structures_path):
        report_parts.append("\n\n---\n\n")
        with open(structures_path, 'r') as f:
            report_parts.append(f.read())
    
    # Add note about sequence diagram
    diagram_path = os.path.join(output_dir, "login_sequence_diagram.puml")
    if os.path.exists(diagram_path):
        report_parts.append("\n\n---\n\n")
        report_parts.append("## Sequence Diagram\n\n")
        report_parts.append("A PlantUML sequence diagram has been generated in `login_sequence_diagram.puml`.\n")
        report_parts.append("To visualize this diagram, use PlantUML or an online tool like https://www.planttext.com/\n")
    
    # Add packet comparison report
    comparison_path = os.path.join(output_dir, "packet_comparison_report.md")
    if os.path.exists(comparison_path):
        report_parts.append("\n\n---\n\n")
        with open(comparison_path, 'r') as f:
            report_parts.append(f.read())
    
    # Write the comprehensive report
    with open(os.path.join(output_dir, "comprehensive_login_analysis.md"), 'w') as f:
        f.write("\n".join(report_parts))
    
    print(f"Comprehensive report generated: comprehensive_login_analysis.md")

def main():
    """Run all analysis scripts and generate a comprehensive report."""
    # Get the directory of this script
    script_dir = os.path.dirname(os.path.abspath(__file__))
    base_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
    analysis_dir = os.path.join(base_dir, "analysis")
    
    # Change to the script directory
    os.chdir(script_dir)
    
    print("Starting Ragnarok Online Login Sequence Analysis")
    print("=" * 50)
    
    # Run the packet analyzer
    run_script("packet_analyzer.py")
    
    # Run the packet structure analyzer
    run_script("packet_structure_analyzer.py")
    
    # Run the sequence diagram generator
    run_script("sequence_diagram_generator.py")
    
    # Run the dump comparison tool with all available dumps
    extracted_data_dir = os.path.join(base_dir, "extracteddata")
    
    # Find all dump JSON files dynamically
    import glob
    json_files = glob.glob(os.path.join(extracted_data_dir, "dump*_packets.json"))
    
    if json_files:
        comparison_output = os.path.join(analysis_dir, "packet_comparison_report.md")
        subprocess.run([
            sys.executable, "compare_dumps.py",
            *json_files,
            "-o", comparison_output
        ])
    print(f"✓ compare_dumps.py completed successfully - report written to {comparison_output}")
    
    # Generate the comprehensive report
    generate_comprehensive_report(analysis_dir)
    
    print("=" * 50)
    print("Analysis complete! All results are available in the following files:")
    print("- analysis/login_sequence_summary.md")
    print("- analysis/packet_structures.md")
    print("- analysis/login_sequence_diagram.puml")
    print("- analysis/packet_comparison_report.md")
    print("- analysis/comprehensive_login_analysis.md (combined report)")

if __name__ == "__main__":
    main()